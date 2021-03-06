// +build sender

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	gcfg "gopkg.in/gcfg.v1"

	"image/png"

	"github.com/golang-mcr/machiavelli/steganogopher"
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

	tmpfile, err := ioutil.TempFile("", "example.jpg")
	if err != nil {
		log.Fatal(err)
	}

	//	defer os.Remove(tmpfile.Name()) // clean up
	log.Printf("temp file: %s\n", tmpfile.Name())
	reader, err := os.Open("steganogopher/_test/terrorcat.png")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	m, err := png.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	err = steganogopher.Encode(tmpfile, m, message, nil)
	if err != nil {
		log.Fatal(err)
	}
	tweet := twitter.Tweet{
		Message: fmt.Sprintf("#%s", GetCurrentHashTag()),
		Image:   tmpfile.Name(),
	}
	fmt.Println("+-+-+-+-+-+-+-+-+-+-+-+-+-+\n|S|t|e|g|a|n|o|G|O|p|h|e|r|\n+-+-+-+-+-+-+-+-+-+-+-+-+-+")
	fmt.Println("Sending encoded message...")
	fmt.Println(tweet.Message)

	err = client.Tweet(tweet)
	if err != nil {
		log.Printf(err.Error())
	}
}
