package main

import (
	"fmt"
	"path"

	"github.com/altid/libs/fs"
	"github.com/dghubble/go-twitter/twitter"
)

func setupStream(s *server, ctrl *fs.Control, tc *twitter.Client) (*twitter.Stream, error) {
	// Convenience Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()

	IDList := getFollows(tc)

	demux.Tweet = func(tweet *twitter.Tweet) {
		if _, ok := IDList[tweet.User.ID]; !ok {
			return
		}

		m, _ := ctrl.MainWriter("main", "feed")
		defer m.Close()

		fmt.Fprintf(m, "%%[@%s](blue) %s [link](#%d)\n", tweet.User.ScreenName, tweet.Text, tweet.ID)
		ctrl.Event(path.Join(*mtpt, *srv, "main"))
	}

	demux.DM = func(dm *twitter.DirectMessage) {

		ctrl.CreateBuffer("@"+dm.Sender.ScreenName, "feed")

		dw, _ := ctrl.MainWriter("@"+dm.Sender.ScreenName, "feed")
		defer dw.Close()

		fmt.Fprintf(dw, "%%[@%s](blue) %s\n", dm.Sender.ScreenName, dm.Text)

		ctrl.Input("@" + dm.Sender.ScreenName)
		ctrl.Event(path.Join(*mtpt, *srv, "@"+dm.Sender.ScreenName))
	}

	// Log to err
	demux.Warning = func(warning *twitter.StallWarning) {
		ew, _ := ctrl.ErrorWriter()

		defer ew.Close()
		fmt.Fprintf(ew, "API error %s: %s\n", warning.Code, warning.Message)
	}

	userParams := &twitter.StreamFilterParams{
		StallWarnings: twitter.Bool(true),
		Follow:        toList(IDList),
	}

	stream, err := tc.Streams.Filter(userParams)
	if err != nil {
		return nil, err
	}

	// Receive messages until stopped or stream quits
	go demux.HandleChan(stream.Messages)
	return stream, nil
}

func getFollows(tc *twitter.Client) map[int64]string {
	friendlist := make(map[int64]string)

	var n int64 = -1
	for n != 0 {
		friends, _, err := tc.Friends.List(&twitter.FriendListParams{
			SkipStatus:          twitter.Bool(true),
			IncludeUserEntities: twitter.Bool(false),
			Count:               200,
			Cursor:              n,
		})

		if err != nil {
			break
		}

		for _, friend := range friends.Users {
			friendlist[friend.ID] = friend.IDStr
		}

		n = friends.NextCursor
	}

	return friendlist
}

func watch(id string, list []string) bool {

	for _, good := range list {
		if id == good {
			return true
		}
	}

	return false
}

func toList(incoming map[int64]string) []string {
	var list []string

	for _, item := range incoming {
		list = append(list, item)
	}

	return list
}
