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
		fmt.Fprintf(m, "%%[@%s](blue) %s #%d\n", tweets[i-1].User.ScreenName, tweets[i-1].Text, tweets[i-1].ID)
	}

	ctrl.Event(path.Join(*mtpt, *srv, "main"))
	return nil
}

func watchFeed(s *server, ctrl *fs.Control, tc *twitter.Client) (*twitter.Stream, error) {
	// Convenience Demux demultiplexed stream messages
	m, err := ctrl.MainWriter("main", "feed")
	if err != nil {
		return nil, err
	}

	defer m.Close()

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		fmt.Fprintf(m, "%%[@%s](blue) %s #%d\n", tweet.User.ScreenName, tweet.Text, tweet.ID)
		ctrl.Event(path.Join(*mtpt, *srv, "main"))
	}

	demux.DM = func(dm *twitter.DirectMessage) {
		ctrl.CreateBuffer(dm.Sender.ScreenName, "feed")
		dw, _ := ctrl.MainWriter(dm.Sender.ScreenName, "feed")

		fmt.Fprintf(dw, "%%[@%s](blue) %s\n", dm.Sender.ScreenName, dm.Text)
		input, _ := fs.NewInput(s, path.Join(*mtpt, *srv), dm.Sender.ScreenName, *debug)
		input.Start()

		dw.Close()
	}

	demux.Event = func(event *twitter.Event) {
		fmt.Fprintf(m, "%%[@%s](green) %s\n", event.Source.ScreenName, event.Event)
	}

	userParams := &twitter.StreamUserParams{
		StallWarnings: twitter.Bool(true),
	}

	stream, err := tc.Streams.User(userParams)
	if err != nil {
		return nil, err
	}

	// Receive messages until stopped or stream quits
	go demux.HandleChan(stream.Messages)
	return stream, nil
}

func tweet(tc *twitter.Client, msg string) {
	tc.Statuses.Update(msg, nil)
}

func retweet(tc *twitter.Client, id int64) {
	tc.Statuses.Retweet(id, nil)
}
