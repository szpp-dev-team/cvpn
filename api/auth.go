// login とか logout とか

package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func (c *Client) Login(username string, password string) error {
	if err := c.login(username, password); err != nil {
		return err
	}

	authParams, err := c.getAuthParams()
	if err != nil {
		return err
	}
	c.authParams = authParams

	if err := saveCookies(c.cookies); err != nil {
		return err
	}

	return nil
}

type SessionError struct{}

func (se *SessionError) Error() string {
	return "Error: Session Error. You have to choose session"
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
		return &SessionError{}
	default: // confirm session
		return errors.New("Oops! You should choose session(todo)")
	}

	return nil
}

func (c *Client) getAuthParams() (map[string][]string, error) {
	doc, err := c.getDoc(
		VpnIndexURL,
		func(req *http.Request, resp *http.Response) error {
			if resp.StatusCode != http.StatusOK {
				if isRedirectedToLogin(resp) {
					return &ErrRedirectedToLogin{NextPath: resp.Header.Get("location"), PrevPath: req.URL.RawPath}
				}
				return errors.New("Error: Status code was not OK")
			}

			return nil
		},
	)
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

func (c *Client) LoadCookiesOrLogin(username, password string) error {
	cookies, err := loadCookies()
	if err != nil || len(cookies) == 0 {
		log.Println("LoadCookiesOrLogin(): Trying login due to no cookie.")
		return c.Login(username, password)
	}
	c.cookies = cookies

	// ファイルから読み込んだクッキーで getAuthParms() が成功したなら return
	authParams, err := c.getAuthParams()
	if err == nil {
		c.authParams = authParams
		log.Println("LoadCookiesOrLogin(): Succeeded getAuthParms() with saved cookie.")
		return nil
	}

	switch err.(type) {
	case *ErrRedirectedToLogin:
		log.Println("LoadCookiesOrLogin(): Trying login due to invalid cookie.")
		return c.Login(username, password)
	default:
		return err
	}
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

	req, _ = http.NewRequest(http.MethodGet, VpnHostRoot+location, nil)
	resp, err = c.request(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to logout")
	}

	_ = deleteCookieFile()
	return nil
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

	fmt.Println("xsauth: " + val)

	return val, nil
}
