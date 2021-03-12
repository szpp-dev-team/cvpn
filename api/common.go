// 共通の api(cookie をつけたリクエストとか)

package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	VpnHostRoot = "https://vpn.inf.shizuoka.ac.jp"
	VpnIndexURL = "https://vpn.inf.shizuoka.ac.jp/dana/home/index.cgi"
)

type Client struct {
	client     *http.Client
	cookies    []string
	authParams map[string][]string
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
	if len(c.cookies) == 0 {
		return nil, errors.New("Error: please login")
	}

	r.Header = map[string][]string{
		"cookie": {strings.Join(c.cookies, "; ")}}

	resp, err := c.client.Do(r)

	return resp, err
}

type ErrRedirectedToLogin struct {
	NextPath string // リダイレクト先のURLパス
	PrevPath string // リダイレクト元のURLパス
}

func (err *ErrRedirectedToLogin) Error() string {
	return fmt.Sprintf("Err: Redirected to login %q (previous: %q)",
		err.NextPath, err.PrevPath)
}

func isRedirectedToLogin(resp *http.Response) bool {
	if resp.StatusCode != http.StatusFound {
		return false
	}
	location, err := resp.Location()
	if err != nil {
		return false
	}
	return strings.Contains(location.String(), "/dana-na/auth/")
}
