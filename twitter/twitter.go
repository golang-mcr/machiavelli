package twitter

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/dghubble/oauth1"
)

type UploadResponse struct {
	MediaId string `json:"media_id_string"`
}

type TwitterStatus struct {
	ID       int    `json:"id"`
	Text     string `json:"text"`
	Lang     string `json:"lang"`
	Entities struct {
		Media []Media `json:"media"`
	} `json:"entities"`
}

type Media struct {
	Id  string `json:"id_str"`
	Url string `json:"media_url"`
}

// Tweet holds the contents of a tweet
type Tweet struct {
	Message string
	Image   string
}

// Client defines the opertations of twitter client
type Client interface {
	Listen(search string) (chan Tweet, func())
	Tweet(tweet Tweet) error
}

type client struct {
	httpClient *http.Client
	config     *Config
}

func (c client) Listen(search string) (chan Tweet, func()) {
	ch := make(chan Tweet)
	var cancel bool
	var lastTweet int
	go func() {
		for !cancel {
			req, _ := http.NewRequest(http.MethodGet, apiURL+apiVersion+timelineURI, nil)
			q := req.URL.Query()
			q.Add("screen_name", search)
			q.Add("count", "1")
			q.Add("include_rts", "false")
			req.URL.RawQuery = q.Encode()

			resp, err := c.httpClient.Do(req)
			//fmt.Println("[+] Polling for tweets")

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

			var tweets []TwitterStatus
			json.Unmarshal(respBody, &tweets)

			for _, tweet := range tweets {
				if tweet.ID > lastTweet {
					lastTweet = tweet.ID
					log.Printf("tweet ID: %d\n", tweet.ID)
					ch <- Tweet{
						Message: tweet.Text,
						Image:   tweet.Entities.Media[0].Id,
					}
				}
			}
		}
	}()

	return ch, func() { cancel = true }
}

func newfileUploadRequest(uri string, paramName string, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, err
}

func (c client) Tweet(tweet Tweet) error {
	if len(tweet.Message) > 140 {
		return errors.New("tweet exceeds 140 character limit")
	}
	url := fmt.Sprintf(apiURL+apiVersion+statusURI+"?status=%s", encodeStatus(&tweet.Message))

	if len(tweet.Image) > 0 {
		req, err := newfileUploadRequest(fmt.Sprintf(mediaApiURL+apiVersion+mediaURI), "media", tweet.Image)
		if err != nil {
			return fmt.Errorf("error building request: %v", err)
		}

		res, err := c.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("%s", res.Status)
		}

		defer res.Body.Close()

		var uploadResponse UploadResponse
		json.NewDecoder(res.Body).Decode(&uploadResponse)
		//fmt.Println(uploadResponse.MediaId)
		//		file, err := ("", uploadResponse.MediaId)
		//		if err != nil {
		//			log.Fatalf("error creating temp file for media id: %s\n", err)
		//		}
		//		src, _ := os.Open(tweet.Image)
		//		io.Copy(file, src)
		os.Rename(tweet.Image, filepath.Join("images", uploadResponse.MediaId))

		url = fmt.Sprintf("%s&media_ids=%s", url, encodeStatus(&uploadResponse.MediaId))
	}

	log.Printf("url: %s\n", url)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("error building request: %v", err)
	}

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
	config1 := oauth1.NewConfig(config.ConsumerKey, config.ConsumerSecret)
	token := oauth1.NewToken(config.AccessToken, config.AccessTokenSecret)
	httpClient = config1.Client(oauth1.NoContext, token)

	return client{httpClient, config}
}
