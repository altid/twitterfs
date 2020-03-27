package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mrjones/oauth"
)

// We need to have our consumer-key and consumer-secret on a server
// Then we can serve up these tokens to users, and they can set
// We stash auth tokens in a db, and return them on lookup
func generateToken() *oauth.AccessToken {
	l := log.New(os.Stderr, "", 0)

	c := oauth.NewConsumer(
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_CONSUMER_SECRET"),
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
			AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
		},
	)

	c.Debug(*debug)

	requestToken, u, err := c.GetRequestTokenAndUrl("oob")
	if err != nil {
		l.Fatal("Unable to complete request. Did you set TWITTER_CONSUMER_KEY and/or TWITTER_CONSUMER_SECRET environment variables?")
	}

	fmt.Println("A link has been created to retreive a 7-Digit key")
	fmt.Printf("\n%s\n", u)
	fmt.Println("Paste your 7-Digit key below, and press enter:")

	var accessToken *oauth.AccessToken
	var count uint

	for {
		if count > 2 {
			l.Fatal("Failed to validate token, please try again later")
		}
		verificationCode := ""
		fmt.Scanln(&verificationCode)

		accessToken, err = c.AuthorizeToken(requestToken, verificationCode)
		if err != nil {
			count++
			l.Println("something went wrong. Please enter your code again")
			continue
		}

		break
	}
	return accessToken
}
