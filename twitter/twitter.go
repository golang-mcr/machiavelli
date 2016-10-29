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

func (c client) Listen(search string) (<-chan Tweet, func()) {
	ch := make(chan Tweet)

	return ch, func() {}
}

func (c client) Tweet(tweet Tweet) error {
	return nil
}

// NewClient takes a http client with oauth credentials
// to make future calls to twitter
func NewClient(httpClient *http.Client) Client {
	return client{httpClient}
}
