package twitter

import "net/http"

// Tweet holds the contents of a tweet
type Tweet struct {
	Message string
}

// Client defines the opertations of twitter client
type Client interface {
	Listen(search string) (<-chan Tweet, func())
	Tweet(tweet Tweet) error
}

type client struct {
	httpClient *http.Client
}

// NewClient takes a http client with oauth credentials
// to make future calls to twitter
func NewClient(httpClient *http.Client) *Client {
}
