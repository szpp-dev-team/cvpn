// login とか logout とか

package api

import (
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func (c *Client) Login(username string, password string) error {
	params := make(url.Values)

	params.Set("tz_offset", "540")
	params.Set("username", username)
	params.Set("password", password)
	params.Set("realm", "Student-Realm")
	params.Set("btnSubmit", "Sign+In")

	if err := c.login(params); err != nil {
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

type SessionError struct {
	requestURL string
}

func (se *SessionError) Error() string {
	return "Error: Session Error. You have to choose session"
}

func (c *Client) login(params url.Values) error {
	const (
		LoginEndpoint = "https://vpn.inf.shizuoka.ac.jp/dana-na/auth/url_3/login.cgi"
		LoginFailed   = "/dana-na/auth/url_3/welcome.cgi?p=failed"
		LoginSucceed  = "/dana/home/index.cgi"
	)

	resp, err := c.client.PostForm(LoginEndpoint, params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c.cookies = getCookies(resp.Header["Set-Cookie"])
	location := resp.Header.Get("location")
	switch location {
	case LoginSucceed:
		return nil
	case LoginFailed:
		return errors.New("Error: Login Failed")
	default:
		return &SessionError{requestURL: location}
	}
}

func (c *Client) getAuthParams() (url.Values, error) {
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

	params := make(url.Values)

	params.Set("xsauth", xsauth)

	return params, nil
}

func (c *Client) LoadCookiesOrLogin(username, password string) error {
	cookies, err := loadCookies()
	if err != nil || len(cookies) == 0 {
		log.Println("LoadCookiesOrLogin(): Trying login due to no cookie.")
		return c.Login(username, password)
	}
	c.cookies = cookies

	// ファイルから読み込んだクッキーで getAuthparams() が成功したなら return
	authParams, err := c.getAuthParams()
	if err == nil {
		c.authParams = authParams
		log.Println("LoadCookiesOrLogin(): Succeeded getAuthparams() with saved cookie.")
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

// if ok == true, continue to login on current device.
// else, stop to login.
func (c *Client) ConfirmSession(ok bool) error {
	const (
		ConfirmEndpoint = "https://vpn.inf.shizuoka.ac.jp/dana-na/auth/url_3/welcome.cgi?p=user%2Dconfirm"
		ContinueLogin   = "セッションを続行します"
		StopLogin       = "キャンセル"
	)

	doc, err := c.getDoc(
		ConfirmEndpoint,
		func(req *http.Request, resp *http.Response) error {
			if resp.StatusCode != http.StatusOK {
				return errors.New("Error: not ok")
			}
			return nil
		},
	)
	if err != nil {
		return err
	}

	formDataStr, err := findFormDataStr(doc)
	if err != nil {
		return err
	}

	params := func() url.Values {
		if ok {
			return url.Values{
				"btnContinue": {ContinueLogin},
			}
		}
		return url.Values{
			"btnCancel": {StopLogin},
		}
	}()
	params["FormDataStr"] = []string{formDataStr}

	if err := c.login(params); err != nil {
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

func findFormDataStr(doc *goquery.Document) (string, error) {
	selection := doc.Find(`
		html >
		body >
		form#DSIDConfirmForm >
		input#DSIDFormDataStr
	`)

	formDataStr, exists := selection.Attr("value")
	if !exists {
		return "", errors.New("Error: FormDataStr not found")
	}

	// fmt.Println("FormDataStr:", formDataStr)

	return formDataStr, nil
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

	// fmt.Println("xsauth: " + val)

	return val, nil
}
