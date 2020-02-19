package main

import (
	"fmt"
	"path"
	"strings"

	cm "github.com/altid/cleanmark"
	"github.com/altid/fslib"
)

var workdir = path.Join(*mtpt, *srv)

type server struct {
}

func (s *server) Open(c *fslib.Control, name string) error {
}

func (s *server) Close(c *fslib.Control, name string) error {
}

func (s *server) Link(c *fslib.Control, from, name string) error {
	return fmt.Errorf("link command not supported, please use open/close\n")
}

func (s *server) Default(c *fslib.Control, cmd, from, m string) error {
}

// input is always sent down raw to the server
func (s *server) Handle(bufname string, l *cm.Lexer) error {
	var m strings.Builder
	for {
		i := l.Next()
		switch i.ItemType {
		case cm.EOF:
			// write m.String() to thing
			return err
		case cm.ErrorText:
		case cm.UrlLink, cm.UrlText, cm.ImagePath, cm.ImageLink, cm.ImageText:
		case cm.ColorText, cm.ColorTextBold:
		case cm.BoldText:
		case cm.EmphasisText:
		case cm.UnderlineText:
		default:
			m.Write(i.Data)
		}
	}
	return fmt.Errorf("Unknown error parsing input encountered")
}
