package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/altid/libs/fs"
	"github.com/altid/libs/markup"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

var workdir = path.Join(*mtpt, *srv)

type server struct {
	tc     *twitter.Client
	handle string
}

// Host consumer-key and consumer-secret on our site
// queries come in with token, secret and we use that to
// oauth validate them to our service, returning the config entry
// create our own client with the auth proxied through our servers
func newServer(token, secret string) *server {
	config := oauth1.NewConfig(
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_CONSUMER_SECRET"),
	)
	acctoken := oauth1.NewToken(token, secret)
	client := config.Client(oauth1.NoContext, acctoken)

	tc := twitter.NewClient(client)
	user, _, _ := tc.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{})

	return &server{tc, "@" + user.ScreenName}
}

// Start a PM
func open(s *server, c *fs.Control, name string) error {
	if name[0] != '@' {
		name = "@" + name
	}

	// Make sure the handle is good, get the ID
	user, _, err := s.tc.Users.Lookup(&twitter.UserLookupParams{
		ScreenName: []string{name[1:]},
	})

	id := user[0].IDStr

	if e := c.CreateBuffer(name, "feed"); e != nil {
		return e
	}

	mw, err := c.MainWriter(name, "feed")
	if err != nil {
		return err
	}

	// Fetch our list of the last 200 DMS
	dms, _, err := s.tc.DirectMessages.EventsList(&twitter.DirectMessageEventsListParams{
		Count: 200,
	})

	// Any DM between us and them, add
	for i := len(dms.Events); i > 0; i-- {
		event := dms.Events[i-1]

		// To us from them
		if event.Message.SenderID == id {
			fmt.Fprintf(mw, "%%[%s](blue) %s\n", name, event.Message.Data.Text)
			continue
		}

		// From us to them
		if event.Message.Target.RecipientID == id {
			fmt.Fprintf(mw, "%%[%s](grey) %s\n", s.handle, event.Message.Data.Text)
			continue
		}
	}

	c.Event(path.Join(*mtpt, *srv, name, "feed"))
	return c.Input(name)
}

func (s *server) Run(c *fs.Control, cmd *fs.Command) error {
	switch cmd.Name {
	case "open":
		// Tweet handles cannot have spaces
		return open(s, c, cmd.Args[0])
	case "close":
		return c.DeleteBuffer(cmd.Args[0], "feed")
	case "tweet":
		msg := strings.Join(cmd.Args, " ")
		_, _, err := s.tc.Statuses.Update(msg, nil)

		return err
	case "rt":
		id, err := strconv.Atoi(cmd.Args[0][1:])
		if err != nil {
			return err
		}

		_, _, err = s.tc.Statuses.Retweet(int64(id), nil)
		return err
	//case "reply" id data...
	//case "love":
	//case "follow":
	//case "msg":

	default:
		return errors.New("Command not supported")
	}
}

func (s *server) Quit() {

}

// Twitter doesn't support any formatting, don't use
func (s *server) Handle(bufname string, l *markup.Lexer) error {
	if bufname[0] == '@' {
		bufname = bufname[1:]
	}

	u, _, err := s.tc.Users.Lookup(
		&twitter.UserLookupParams{
			ScreenName: []string{bufname},
		},
	)
	if err != nil {
		return err
	}

	input, err := l.String()
	if err != nil {
		return err
	}

	return dm(s.tc, u[0].IDStr, input)
}
