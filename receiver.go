// +build receiver

package main

import (
	"flag"
	"fmt"
	"os"
	"net/http"

	"github.com/golang-mcr/machiavelli/twitter"
	gcfg "gopkg.in/gcfg.v1"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "", "configuration file (.gcfg)")
	flag.Parse()

	var cfg Config
	if err := gcfg.ReadFileInto(&cfg, configFile); err != nil {
		fmt.Fprintf(os.Stderr, "error getting config variables: %v", err)
		return
	}

		client := twitter.NewClient(http.DefaultClient, &cfg.Twitter)
	_, stop := client.Listen("test")
	stop()
	fmt.Println("machiavelli")
}
