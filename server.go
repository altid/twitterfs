package main

import (
	"context"
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
	cancel context.CancelFunc
	tc     *twitter.Client
	handle string
}

// Host consumer-key and consumer-secret on our site
// queries come in with token, secret and we use that to
// oauth validate them to our service, returning the config entry
// create our own client with the auth proxied through our servers
func newServer(cancel context.CancelFunc, token, secret, handle string) *server {
	config := oauth1.NewConfig(
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_CONSUMER_SECRET"),
	)
	acctoken := oauth1.NewToken(token, secret)
	client := config.Client(oauth1.NoContext, acctoken)

	tc := twitter.NewClient(client)
	return &server{cancel, tc, handle}
}

// Start a PM
func (s *server) Open(c *fs.Control, name string) error {
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

	input, err := fs.NewInput(s, path.Join(*mtpt, *srv), name, *debug)
	if err != nil {
		return err
	}

	input.Start()
	return nil
}

func (s *server) Close(c *fs.Control, name string) error {
	return c.DeleteBuffer(name, "feed")
}

func (s *server) Link(c *fs.Control, from, name string) error {
	return errors.New("link command not supported, please use open/close")
}

func (s *server) Default(c *fs.Control, cmd *fs.Command) error {
	switch cmd.Name {
	case "tweet":
		msg := strings.Join(cmd.Args, " ")
		s.tc.Statuses.Update(msg, nil)
	case "rt":
		id, err := strconv.Atoi(cmd.Args[0][1:])
		if err != nil {
			return err
		}

		s.tc.Statuses.Retweet(int64(id), nil)
		//case "reply" id data...
		//case "love":
		//case "follow":
		//case "msg":
	}
	return nil
}

func (s *server) Refresh(c *fs.Control) error {
	return nil
}

func (s *server) Restart(c *fs.Control) error {
	return nil
}

func (s *server) Quit() {
	s.cancel()
}

// Twitter doesn't support any formatting, don't use
func (s *server) Handle(bufname string, l *markup.Lexer) error {
	var m strings.Builder

	for {
		i := l.Next()
		switch i.ItemType {
		case markup.EOF:
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

			return dm(s.tc, u[0].IDStr, m.String())
		case markup.ErrorText:
		case markup.URLLink, markup.URLText, markup.ImagePath, markup.ImageLink, markup.ImageText:
		case markup.ColorText, markup.ColorTextBold:
		case markup.BoldText:
		case markup.EmphasisText:
		case markup.UnderlineText:
		default:
			m.Write(i.Data)
		}
	}
}
