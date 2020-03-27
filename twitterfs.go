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
)

func main() {
	flag.Parse()
	if flag.Lookup("h") != nil {
		flag.Usage()
		os.Exit(1)
	}

	conf := &struct {
		ListenAddress types.ListenAddress
		Logdir        types.Logdir
		Handle        string `Enter your current Twitter handle (@foo)`
		Token         string
		Secret        string
	}{"none", "none", "", "none", "none"}

	if *setup {
		at := generateToken()
		conf.Token = at.Token
		conf.Secret = at.Secret
		if e := config.Create(conf, *srv, "", *debug); e != nil {
			log.Fatal(e)
		}

		os.Exit(0)
	}

	if e := config.Marshal(conf, *srv, "", *debug); e != nil {
		log.Fatal(e)
	}

	if conf.Token == "none" || conf.Secret == "none" {
		at := generateToken()
		
		conf.Token = at.Token
		conf.Secret = at.Secret
		
		log.Printf("To skip this step, run %s -conf and store keys in a conf", os.Args[0])
	}

	ctx, cancel := context.WithCancel(context.Background())
	s := newServer(cancel, conf.Token, conf.Secret, conf.Handle)

	ctrl, err := fs.CreateCtlFile(ctx, s, string(conf.Logdir), *mtpt, *srv, "feed", *debug)
	if err != nil {
		log.Fatal(err)
	}

	defer ctrl.Cleanup()
	ctrl.SetCommands(TwitterCommands...)

	if e := setupMain(ctrl, s.tc); e != nil {
		log.Fatal(e)
	}

	stream, err := setupStream(s, ctrl, s.tc)
	if err != nil {
		log.Fatal(err)
	}

	defer stream.Stop()

	if e := ctrl.Listen(); e != nil {
		log.Fatal(e)
	}
}
