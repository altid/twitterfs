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
func generateToken() {
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
		log.Fatal(err)
	}
	fmt.Printf("Please go to %s and retreive the 7-digit key\n", u)
	fmt.Println("Enter your key below, and press enter:")

	var accessToken *oauth.AccessToken
	var count uint

	for {
		if count > 2 {
			fmt.Println("Failed to validate token, please try again later")
			os.Exit(0)
		}
		verificationCode := ""
		fmt.Scanln(&verificationCode)

		accessToken, err = c.AuthorizeToken(requestToken, verificationCode)
		if err != nil {
			count++
			fmt.Println("something went wrong. Please enter your code again")
			continue
		}

		break
	}
	fmt.Printf("Success! Please add this value to your altid/config:\n\nserver=%s token=%s auth=password password=%s\n",
		*srv,
		accessToken.Token,
		accessToken.Secret,
	)
	os.Exit(0)
}
