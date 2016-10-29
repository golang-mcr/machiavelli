package twitter

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
  "log"
)

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
	config     *Config
}

func (c client) Listen(search string) (<-chan Tweet, func()) {
	ch := make(chan Tweet)
	var cancel bool
	for !cancel {
		req, _ := http.NewRequest(http.MethodGet, apiURL+apiVersion+timelineURI, nil)
		q := req.URL.Query()
		q.Add("screen_name", search)
		q.Add("count", "1")
		q.Add("include_rts", "false")
		req.URL.RawQuery = q.Encode()

		oa := NewOAuthDetails(c.config, "something")
		req.Header.Set(authHeader, fmt.Sprintf("%s", oa))

		resp, err := c.httpClient.Do(req)
		fmt.Println("[+] Polling for tweets")

		if err != nil {
			log.Printf("%v", err)
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf(resp.Status)
		}

		defer resp.Body.Close()
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("%v", err)
		}

		var tweet = Tweet{Message: string(respBody)}

		ch <- tweet
	}

	return ch, func() { cancel = true }
}

func (c client) Tweet(tweet Tweet) error {
	if len(tweet.Message) > 140 {
		return errors.New("tweet exceeds 140 character limit")
	}

	oa := NewOAuthDetails(c.config, tweet.Message)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(apiURL+apiVersion+statusURI+"?status=%s", encodeStatus(&tweet.Message)), nil)
	if err != nil {
		return fmt.Errorf("error building request: %v", err)
	}
	req.Header.Set(authHeader, fmt.Sprintf("%s", oa))

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", res.Status)
	}

	return nil
}

// NewClient takes a http client with oauth credentials
// to make future calls to twitter
func NewClient(httpClient *http.Client, config *Config) Client {
	return client{httpClient, config}
}
