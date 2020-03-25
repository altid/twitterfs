package main

import (
	"time"

	"github.com/dghubble/go-twitter/twitter"
)

func dm(tc *twitter.Client, id, msg string) error {
	_, _, err := tc.DirectMessages.EventsNew(&twitter.DirectMessageEventsNewParams{
		Event: &twitter.DirectMessageEvent{
			Type:      "message_create",
			CreatedAt: time.Now().String(),
			Message: &twitter.DirectMessageEventMessage{
				Data: &twitter.DirectMessageData{
					Text: msg,
				},
				Target: &twitter.DirectMessageTarget{
					RecipientID: id,
				},
			},
		},
	})

	return err
}

func tm(tc *twitter.Client, id int64, msg string) error {
	_, _, err := tc.Statuses.Update(msg, &twitter.StatusUpdateParams{
		InReplyToStatusID: id,
	})

	return err
}
