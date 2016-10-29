// +build sender

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	gcfg "gopkg.in/gcfg.v1"

	"github.com/golang-mcr/machiavelli/twitter"
)

func main() {
	var configFile, message string
	flag.StringVar(&configFile, "config", "", "configuration file (.gcfg)")
	flag.StringVar(&message, "message", "", "message to be sent")
	flag.Parse()

	var cfg Config
	if err := gcfg.ReadFileInto(&cfg, configFile); err != nil {
		fmt.Fprintf(os.Stderr, "error getting config variables: %v from %v\n", err, configFile)
		return
	}

	client := twitter.NewClient(http.DefaultClient, &cfg.Twitter)

	tweet := twitter.Tweet{
		Message: message + " #" + GetCurrentHashTag(),
	}
	fmt.Println("+-+-+-+-+-+-+-+-+-+-+-+-+-+\n|S|t|e|g|a|n|o|G|O|p|h|e|r|\n+-+-+-+-+-+-+-+-+-+-+-+-+-+")
	fmt.Println("Sending encoded message...")
	fmt.Println(tweet.Message)

	err := client.Tweet(tweet)
	if err != nil {
		log.Printf(err.Error())
	}
}
