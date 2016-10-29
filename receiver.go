// +build receiver

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
	var configFile string
	flag.StringVar(&configFile, "config", "", "configuration file (.gcfg)")
	flag.Parse()

	var cfg config
	if err := gcfg.ReadFileInto(&cfg, configFile); err != nil {
		fmt.Fprintf(os.Stderr, "error getting config variables: %v", err)
		return
	}

	//auth := &oauth2.Config{}
	// token := &oauth2.Token{AccessToken: cfg.twitter.accessToken}

	config := oauth1.NewConfig(cfg.twitter.consumerKey, cfg.twitter.consumerSecret)
	token := oauth1.NewToken(cfg.twitter.accessToken, cfg.twitter.accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// OAuth2 http.Client will automatically authorize Requests
	httpClient := auth.Client(oauth2.NoContext, token)
	client := twitter.NewClient(httpClient)
	_, stop := client.Listen("test")
	stop()
	fmt.Println("machiavelli")
}
