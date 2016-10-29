package twitter

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func getmsg() {

	client := &http.Client{}

	req, _ := http.NewRequest("GET", "https://api.twitter.com/1.1/statuses/user_timeline.json", nil)
	q := req.URL.Query()
	q.Add("screen_name", "gazj")
	q.Add("count", "1")
	q.Add("include_rts", "false")
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer XXXXXXXXXXXXXXXXXXXXXXXXXXX")

	resp, err := client.Do(req)

	if err != nil {
		// hurrrr
	}

	defer resp.Body.Close()
	resp_body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))

}
