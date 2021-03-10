// 共通の api(cookie をつけたリクエストとか)

package api

import (
	"net/http"
	"strings"
)

const VpnHostRoot = "https://vpn.inf.shizuoka.ac.jp"

type Client struct {
	client  *http.Client
	cookies []string
}

func NewClient() *Client {
	client := new(http.Client)

	// redirect をしない
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &Client{
		client: client,
	}
}

// need cookies
func (c *Client) request(r *http.Request) (*http.Response, error) {
	r.Header = map[string][]string{
		"cookie": {strings.Join(c.cookies, "; ")}}

	resp, err := c.client.Do(r)

	return resp, err
}
