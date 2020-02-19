package main

import (
	"errors"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/altid/libs/fs"
	"github.com/altid/libs/markup"
)

var workdir = path.Join(*mtpt, *srv)

type server struct {
}

func (s *server) Open(c *fs.Control, name string) error {
	return nil
}

func (s *server) Close(c *fs.Control, name string) error {
	return nil
}

func (s *server) Link(c *fs.Control, from, name string) error {
	return errors.New("link command not supported, please use open/close")
}

func (s *server) Default(c *fs.Control, cmd, from, m string) error {
	return nil
}

// input is always sent down raw to the server
func (s *server) Handle(bufname string, l *markup.Lexer) error {
	var m strings.Builder
	for {
		i := l.Next()
		switch i.ItemType {
		case markup.EOF:
			// write m.String() to thing
			return io.EOF
		case markup.ErrorText:
		case markup.UrlLink, markup.UrlText, markup.ImagePath, markup.ImageLink, markup.ImageText:
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
