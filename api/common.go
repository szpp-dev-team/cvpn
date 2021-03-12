// 共通の api(cookie をつけたリクエストとか)

package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

	r.Header.Set("cookie", strings.Join(c.cookies, "; "))
	r.Form = mergeMap(r.Form, c.authParams)

	return c.client.Do(r)
}

// get response body and return *goquery.Document
func (c *Client) getDoc(url string, f func(req *http.Request, resp *http.Response) error) (*goquery.Document, error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	resp, err := c.request(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if f != nil {
		if err := f(req, resp); err != nil {
			return nil, err
		}
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func mergeMap(org, tgt map[string][]string) map[string][]string {
	newmap := make(map[string][]string)

	for k, v := range org {
		newmap[k] = v
	}
	for k, v := range tgt {
		newmap[k] = append(newmap[k], v...)
	}

	return newmap
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
