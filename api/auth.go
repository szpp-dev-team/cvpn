// login とか logout とか

package api

import (
	"errors"
	"net/http"
	"strings"
)

func (self *Client) Login(username string, password string) error {
	const (
		LoginPayload = "https://vpn.inf.shizuoka.ac.jp/dana-na/auth/url_3/login.cgi"
		LoginFailed  = "/dana-na/auth/url_3/welcome.cgi?p=failed"
		LoginSucceed = "/dana/home/index.cgi"
	)

	parms := map[string][]string{
		"tz_offset": {"540"},
		"username":  {username},
		"password":  {password},
		"realm":     {"Student-Realm"},
		"btnSubmit": {"Sign+In"},
	}

	resp, err := self.client.PostForm(LoginPayload, parms)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusFound {
		return errors.New("Error: Please logout and login")
	}

	location := resp.Header.Get("location")
	switch location {
	case LoginSucceed:
		// Header は map[string][]string の拡張
		self.cookies = getCookies(resp.Header["Set-Cookie"])
	case LoginFailed:
		return errors.New("Error: Login Failed")
	default: // confirm session
		return errors.New("Oops! You should choose session(todo)")
	}

	return nil
}

func (self *Client) Logout() error {
	const LogoutPayload = "https://vpn.inf.shizuoka.ac.jp/dana-na/auth/logout.cgi"

	req, err := http.NewRequest(http.MethodGet, LogoutPayload, nil)
	if err != nil {
		return err
	}

	resp, err := self.request(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusFound {
		return errors.New("failed to logout")
	}

	self.cookies = getCookies(resp.Header["Set-Cookie"])
	location := resp.Header.Get("location")

	req, err = http.NewRequest(http.MethodGet, VpnHostRoot+location, nil)
	if err != nil {
		return err
	}

	resp, err = self.request(req)
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
