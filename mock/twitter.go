package mock

import "github.com/golang-mcr/machiavelli/twitter"

// Client mocks the twitter client
type Client struct {
	ListenFunc func(string) (<-chan twitter.Tweet, func())
	TweetFunc  func(twitter.Tweet) error
}

func (c *Client) Listen(search string) (<-chan twitter.Tweet, func()) {
	return c.ListenFunc(search)
}

func (c *Client) Tweet(tweet twitter.Tweet) error {
	return c.TweetFunc(tweet)
}
