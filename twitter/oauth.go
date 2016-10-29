package twitter

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	apiURL          string = "https://api.twitter.com/"
	apiVersion      string = "1.1"
	statusURI       string = "/statuses/update.json"
	timelineURI     string = "/statuses/user_timeline.json"
	searchURI       string = "/search/tweets.json"
	signatureMethod string = "HMAC-SHA1"
	oAuthVersion    string = "1.0"
	authHeader      string = "Authorization"
)

// OAuthDetails contains a valid set of OAuth details based on credentials from a config file.
type OAuthDetails struct {
	ConsumerKey,
	Nonce,
	Signature,
	SignatureMethod,
	Timestamp,
	Token,
	Version string
}

func (oa OAuthDetails) String() string {
	return fmt.Sprintf(`OAuth oauth_consumer_key="%s", oauth_nonce="%s", oauth_signature="%s", oauth_signature_method="%s", oauth_timestamp="%s", oauth_token="%s", oauth_version="%s"`,
		oa.ConsumerKey,
		oa.Nonce,
		oa.Signature,
		signatureMethod,
		oa.Timestamp,
		oa.Token,
		oAuthVersion,
	)
}

// NewOAuthDetails collects a valid set of OAuth details based on credentials passed from a config file.
func NewOAuthDetails(c *Config, status string) *OAuthDetails {
	oa := new(OAuthDetails)
	oa.ConsumerKey = c.ConsumerKey
	oa.Nonce = *generateNonce()
	oa.SignatureMethod = signatureMethod
	oa.Timestamp = *generateTimestamp()
	oa.Token = c.AccessToken
	oa.Version = oAuthVersion
	oa.generateSignature(status, c)

	return oa
}

func generateNonce() *string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	s := string(b)
	return &s
}

func generateTimestamp() *string {
	t := time.Now().Unix()
	f := fmt.Sprintf("%d", t)
	return &f
}

func (oa *OAuthDetails) generateSignature(status string, config *Config) {
	params := collectParams(oa, encodeStatus(&status))
	baseString := generateBaseString(params)
	sig := sign(baseString, config)
	oa.Signature = url.QueryEscape(*sig)
}

func collectParams(oa *OAuthDetails, status string) *string {
	params := fmt.Sprintf(`oauth_consumer_key=%s&oauth_nonce=%s&oauth_signature_method=%s&oauth_timestamp=%s&oauth_token=%s&oauth_version=%s&status=%s`,
		oa.ConsumerKey,
		oa.Nonce,
		signatureMethod,
		oa.Timestamp,
		oa.Token,
		oAuthVersion,
		status,
	)
	return &params
}

func generateBaseString(params *string) *string {
	var b, u bytes.Buffer
	b.WriteString(http.MethodPost)
	b.WriteString("&")

	u.WriteString(apiURL)
	u.WriteString(apiVersion)
	u.WriteString(statusURI)

	b.WriteString(url.QueryEscape(u.String()))
	b.WriteString("&")
	b.WriteString(encodeStatus(params))
	baseString := b.String()

	return &baseString
}

func generateSigningKey(c *Config) *string {
	var b bytes.Buffer
	b.WriteString(url.QueryEscape(c.ConsumerSecret))
	b.WriteString("&")
	b.WriteString(url.QueryEscape(c.AccessTokenSecret))
	key := b.String()

	return &key
}

func sign(baseString *string, c *Config) *string {
	signingKey := generateSigningKey(c)
	h := hmac.New(sha1.New, []byte(*signingKey))
	h.Write([]byte(*baseString))
	result := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return &result
}

func encodeStatus(status *string) string {
	return strings.Replace(url.QueryEscape(*status), "+", "%20", -1)
}
