package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/altid/libs/config"
	"github.com/altid/libs/config/types"
	"github.com/altid/libs/fs"
)

var (
	mtpt  = flag.String("p", "/tmp/altid", "Path for filesystem")
	srv   = flag.String("s", "twitter", "Name of service")
	debug = flag.Bool("d", false, "enable debug logging")
	setup = flag.Bool("conf", false, "Run configuration setup")
	token = flag.Bool("t", false, "Fetch an oauth token")
)

func main() {
	flag.Parse()
	if flag.Lookup("h") != nil {
		flag.Usage()
		os.Exit(1)
	}

	conf := &struct {
		Logdir types.Logdir
		Auth   types.Auth
		User   string
		Token  string
	}{"none", "password", "none", "none"}

	if *setup {
		if e := config.Create(conf, *srv, "", *debug); e != nil {
			log.Fatal(e)
		}

		os.Exit(0)
	}

	if e := config.Marshal(conf, *srv, "", *debug); e != nil {
		log.Fatal(e)
	}

	ctx, cancel := context.WithCancel(context.Background())
	s := newServer(cancel, conf.Token, string(conf.Auth))

	ctrl, err := fs.CreateCtlFile(ctx, s, string(conf.Logdir), *mtpt, *srv, "feed", *debug)
	if err != nil {
		log.Fatal(err)
	}

	defer ctrl.Cleanup()
	ctrl.SetCommands(TwitterCommands...)

	if e := setupMain(ctrl, s.tc); e != nil {
		log.Fatal(e)
	}

	stream, err := watchFeed(s, ctrl, s.tc)
	if err != nil {
		log.Fatal(err)
	}

	defer stream.Stop()

	if e := ctrl.Listen(); e != nil {
		log.Fatal(e)
	}
}
