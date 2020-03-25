package main

import (
	"fmt"
	"path"

	"github.com/altid/libs/fs"
	"github.com/dghubble/go-twitter/twitter"
)

func setupMain(ctrl *fs.Control, tc *twitter.Client) error {
	ctrl.CreateBuffer("main", "feed")

	m, err := ctrl.MainWriter("main", "feed")
	if err != nil {
		return err
	}

	defer m.Close()

	tweets, _, err := tc.Timelines.HomeTimeline(&twitter.HomeTimelineParams{Count: 50})
	if err != nil {
		return err
	}

	// They come in backwards
	for i := len(tweets); i > 0; i-- {
		fmt.Fprintf(m, "%%[@%s](blue) %s [link](#%d)\n", tweets[i-1].User.ScreenName, tweets[i-1].Text, tweets[i-1].ID)
	}

	ctrl.Event(path.Join(*mtpt, *srv, "main"))
	return nil
}
