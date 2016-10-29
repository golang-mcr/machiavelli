package twitter

import (
	"errors"
	"fmt"
	"net/http"
	"io/ioutil"
)

// Tweet holds the contents of a tweet
type Tweet struct {
	Message string
}

// Client defines the opertations of twitter client
type Client interface {
	Listen(search string) (<-chan Tweet, func())
	Tweet(tweet string) error
}

// NewClient takes a http client with oauth credentials
// to make future calls to twitter
// func NewClient(httpClient *http.Client) Client {
// 	return client{httpClient}
// }

type client struct {
	httpClient *http.Client
	config     *Config
}

func (c client) Listen(search string) (<-chan Tweet, func()) {
	return make(chan Tweet), func() { return }
	ch := make(chan Tweet)
	var cancel bool
	for (!cancel){
			req, _ := http.NewRequest(http.MethodGet, apiURL+apiVersion+timelineURI, nil)
			q := req.URL.Query()
			q.Add("screen_name", "hackmanchester")
			q.Add("count", "1")
			q.Add("include_rts", "false")
			req.URL.RawQuery = q.Encode()

			oa := NewOAuthDetails(c.config, "something")
			req.Header.Set(authHeader, fmt.Sprintf("%s", oa))

			resp, _ := c.httpClient.Do(req)
			// if err != nil {
			// 	return fmt.Errorf("%v", err)
			// }
			// if res.StatusCode != http.StatusOK {
			// 	return fmt.Errorf("%s", res.Status)
			// }

			defer resp.Body.Close()
			resp_body, _ := ioutil.ReadAll(resp.Body)
			// if err != nil {
			// 	return fmt.Errorf("%v", err)
			// }

			var tweet = Tweet{Message: string(resp_body)}

			ch <- tweet
	}

	return ch, func() { cancel = true }
}

func (c client) Tweet(tweet string) error {
	if len(tweet) > 140 {
		return errors.New("tweet exceeds 140 character limit")
	}

	oa := NewOAuthDetails(c.config, tweet)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(apiURL+apiVersion+statusURI+"?status=%s", encodeStatus(&tweet)), nil)
	if err != nil {
		return fmt.Errorf("error building request: %v", err)
	}
	req.Header.Set(authHeader, fmt.Sprintf("%s", oa))

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	defer res.Body.Close()

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
