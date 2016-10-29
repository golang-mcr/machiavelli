package twitter

// Tweet holds the contents of a tweet
type Tweet struct {
}

// Client defines the opertations of twitter client
type Client interface {
	Listen(search string) (<-chan Tweet, func())
	Tweet(tweet Tweet) error
}
