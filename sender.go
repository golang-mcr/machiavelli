// +build sender

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang-mcr/machiavelli/twitter"

	"golang.org/x/oauth2"
	gcfg "gopkg.in/gcfg.v1"
)

func main() {
	var configFile, message string
	flag.StringVar(&message, "message", "", "message to be sent")
	flag.Parse()

	client := &mock.Client{
		TweetFunc: func(tweet twitter.Tweet) error {
			fmt.Printf("tweeted: %s\n", tweet.Message)
			return nil
		}
	}

	tweet := twitter.Tweet{
		Message: message,
	}
	client.Tweet(tweet)

	fmt.Println("machiavelli")
}
