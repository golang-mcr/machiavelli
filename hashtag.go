package main

import (
    "math/rand"
    "time"
)

func GetCurrentHashTag() string {
    hashtags := []string{
      "Nulla",
      "venenatis",
      "viverra",
      "pretium",
      "Proin",
      "luctus",
      "ornare",
      "dapibus",
    }

    // TODO: allow seed to be configured per client
    rand.Seed(31337 + time.Now().Unix())

    index := rand.Intn(len(hashtags))

    return hashtags[index]
}
