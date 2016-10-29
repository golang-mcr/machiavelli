// +build sender

package main

import (
	"flag"
	"fmt"

	"github.com/golang-mcr/machiavelli/mock"
	"github.com/golang-mcr/machiavelli/twitter"
)

func main() {
	var message string
	flag.StringVar(&message, "message", "", "message to be sent")
	flag.Parse()

	client := &mock.Client{
		TweetFunc: func(tweet twitter.Tweet) error {
			fmt.Printf("tweeted: %s\n", tweet.Message)
			return nil
		},
	}

	tweet := twitter.Tweet{
		Message: message,
	}
	client.Tweet(tweet)

	fmt.Println("machiavelli")
}
