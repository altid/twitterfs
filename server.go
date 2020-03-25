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
}

// Host consumer-key and consumer-secret on our site
// queries come in with token, secret and we use that to
// oauth validate them to our service, returning the config entry
// create our own client with the auth proxied through our servers
func newServer(cancel context.CancelFunc, token, secret string) *server {
	config := oauth1.NewConfig(
		os.Getenv("TWITTER-CONSUMER-KEY"),
		os.Getenv("TWITTER-CONSUMER-SECRET"),
	)
	acctoken := oauth1.NewToken(token, secret)
	client := config.Client(oauth1.NoContext, acctoken)

	tc := twitter.NewClient(client)
	return &server{cancel, tc}
}

// Start a PM
func (s *server) Open(c *fs.Control, name string) error {
	if name[0] == '@' {
		name = name[1:]
	}

	if e := c.CreateBuffer(name, "feed"); e != nil {
		return e
	}

	input, err := fs.NewInput(s, path.Join(*mtpt, *srv), name, *debug)
	if err != nil {
		return err
	}

	input.Start()
	return nil
}

func (s *server) Close(c *fs.Control, name string) error {
	return nil
}

func (s *server) Link(c *fs.Control, from, name string) error {
	return errors.New("link command not supported, please use open/close")
}

func (s *server) Default(c *fs.Control, cmd *fs.Command) error {
	switch cmd.Name {
	case "tweet":
		tweet(s.tc, strings.Join(cmd.Args, " "))
	case "rt":
		id, err := strconv.Atoi(cmd.Args[0][1:])
		if err != nil {
			return err
		}

		retweet(s.tc, int64(id))
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

// input is always sent down raw to the server
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
	return fmt.Errorf("Unknown error parsing input encountered")
}
