package main

import (
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
  "log"
)

func GetCurrentHashTag() string {
	content, err := ioutil.ReadFile("diceware.lst")
	if err != nil {
		log.Printf("%v", err)
	}
	hashtags := strings.Split(string(content), "\n")

	// TODO: allow seed to be configured per client
	rand.Seed(31337 + time.Now().Unix())

	index := rand.Intn(len(hashtags))

	return hashtags[index]
}
