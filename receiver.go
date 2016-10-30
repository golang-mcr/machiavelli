// +build receiver

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/golang-mcr/machiavelli/steganogopher"
	"github.com/golang-mcr/machiavelli/twitter"
	gcfg "gopkg.in/gcfg.v1"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "", "configuration file (.gcfg)")
	flag.Parse()

	var cfg Config
	if err := gcfg.ReadFileInto(&cfg, configFile); err != nil {
		fmt.Fprintf(os.Stderr, "error getting config variables: %v\n", err)
		return
	}

	fmt.Println("+-+-+-+-+-+-+-+-+-+-+-+-+-+\n|S|t|e|g|a|n|o|G|O|p|h|e|r|\n+-+-+-+-+-+-+-+-+-+-+-+-+-+")
	fmt.Println("Listening for tweets...")

	client := twitter.NewClient(http.DefaultClient, &cfg.Twitter)
	tweets, stop := client.Listen("@go_machiavelli")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for {
		select {
		case tweet := <-tweets:
			log.Printf("tweet: %s %s\n", tweet.Message, tweet.Image)

			tempFile, err := os.Open(filepath.Join("images", tweet.Image))
			if err != nil {
				log.Printf("error opening image: %v\n", err)
				continue
			}
			_, message, err := steganogopher.DecodeAndRead(tempFile)
			if err != nil {
				//log.Printf("error decoding image: %v, file: %s\n", err, tempFile.Name())
				log.Printf("error decoding image: %v\n", err)
				continue
			}
			log.Println("------------- message -------------")
			log.Println(message)
		case <-c:
			stop()
			return
		}
	}
}
