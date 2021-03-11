// login とか logout とか

package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (c *Client) Login(username string, password string) error {
	if err := c.login(username, password); err != nil {
		return err
	}

	_, err := c.getAuthParms()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) login(username, password string) error {
	const (
		LoginEndpoint = "https://vpn.inf.shizuoka.ac.jp/dana-na/auth/url_3/login.cgi"
		LoginFailed   = "/dana-na/auth/url_3/welcome.cgi?p=failed"
		LoginSucceed  = "/dana/home/index.cgi"
	)

	parms := map[string][]string{
		"tz_offset": {"540"},
		"username":  {username},
		"password":  {password},
		"realm":     {"Student-Realm"},
		"btnSubmit": {"Sign+In"},
	}

	resp, err := c.client.PostForm(LoginEndpoint, parms)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusFound {
		return errors.New("Error: Please logout and login")
	}

	location := resp.Header.Get("location")
	switch location {
	case LoginSucceed:
		c.cookies = getCookies(resp.Header["Set-Cookie"])
	case LoginFailed:
		return errors.New("Error: Login Failed")
	default: // confirm session
		return errors.New("Oops! You should choose session(todo)")
	}

	return nil
}

func (c *Client) getAuthParms() (map[string][]string, error) {
	req, err := http.NewRequest(http.MethodGet, VpnIndexURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.request(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Error: Status code was not OK")
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	xsauth, err := findXsauth(doc)
	if err != nil {
		return nil, err
	}

	params := map[string][]string{
		"xsauth": {xsauth},
	}

	return params, nil
}

func (c *Client) Logout() error {
	const LogoutEndpoint = "https://vpn.inf.shizuoka.ac.jp/dana-na/auth/logout.cgi"

	req, err := http.NewRequest(http.MethodGet, LogoutEndpoint, nil)
	if err != nil {
		return err
	}

	resp, err := c.request(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusFound {
		return errors.New("failed to logout")
	}

	c.cookies = getCookies(resp.Header["Set-Cookie"])
	location := resp.Header.Get("location")

	req, err = http.NewRequest(http.MethodGet, VpnHostRoot+location, nil)
	if err != nil {
		return err
	}

	resp, err = c.request(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to logout")
	}

	return nil
}

func getCookies(cookies []string) []string {
	var parsedCookies []string
	for _, line := range cookies {
		tokens := strings.Split(line, "; ")
		parsedCookies = append(parsedCookies, tokens[0])
	}

	return parsedCookies
}

func findXsauth(doc *goquery.Document) (string, error) {
	selection := doc.Find(`
		html >
		body >
		table#table_useruimenu_10.tdContent >
		tbody >
		tr >
		td >
		form#expandForm >
		input#xsauth_395
	`)

	val, exists := selection.Attr("value")
	if !exists {
		return "", errors.New("Error: xsauth not found")
	}

	fmt.Println(val)

	return val, nil
}
