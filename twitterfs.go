package main

import (
	"flag"
	"log"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/altid/fslib"
	"golang.org/x/oauth2"
)

var (
	mtpt = flag.String("p", "/tmp/altid", "Path for filesystem")
	srv  = flag.String("s", "twitter", "Name of service")
)

func main() {
	flag.Parse()
	if flag.Lookup("h") != nil {
		flag.Usage()
		os.Exit(1)
	}

	config, err := newConfig()
	if err != nil {
		log.Fatal(err)
	}

	oa := &oauth2.Config{}
	token := &oauth2.Token{
		AccessToken: config.token,
	}
	httpClient := oa.Client(oauth2.NoContext, token)
	tc := twitter.NewClient(httpClient)
	s := &server{}
	if err != nil {
		log.Fatal("Error initiating Twitter session %v", err)
	}
	ctrl, err := fslib.CreateCtrlFile(s, config.log, *mtpt, *srv, "feed")
	defer ctrl.Cleanup()
	if err != nil {
		log.Fatal(err)
	}
	ctrl.CreateBuffer("main", "feed")
	ctrl.Listen()
}
