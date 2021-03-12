package api

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const UploadEndpoint = "https://vpn.inf.shizuoka.ac.jp/dana/fb/smb/wu.cgi"

// dst は report をルートとしたパス
func (c *Client) UploadFile(src, dst string) error {
	doc, err := c.getDoc(UploadEndpoint, nil)
	if err != nil {
		return err
	}

	trackid, err := findTrackid(doc)
	if err != nil {
		return err
	}

	return c.uploadFile(src, dst, trackid)
}

func (c *Client) uploadFile(src, dst, trackid string) error {
	contentType, body, err := createFormBody(src)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, UploadEndpoint, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	req.Form = map[string][]string{
		"dir":     {url.QueryEscape(strings.Replace("report"+dst, "/", "\\", -1))},
		"acttype": {"upload"},
		"confirm": {"yes"},
		"trackid": {trackid},
	}

	resp, err := c.request(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error: statuscode was not OK but %d", resp.StatusCode)
	}

	/*
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		fmt.Println(string(b))
	*/

	return nil
}

// https://medium.com/eureka-engineering/multipart-file-upload-in-golang-c4a8eb15a3ee
// return {Content-Type}, body, error
func createFormBody(path string) (string, *bytes.Buffer, error) {
	fp, err := os.Open(path)
	if err != nil {
		return "", nil, err
	}
	defer fp.Close()

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	filename := filepath.Base(path)
	fw, err := w.CreateFormFile("file1", filename)
	if err != nil {
		return "", nil, err
	}

	if _, err := io.Copy(fw, fp); err != nil {
		return "", nil, err
	}

	contentType := w.FormDataContentType()

	return contentType, body, w.Close()
}

func findTrackid(doc *goquery.Document) (string, error) {
	selection := doc.Find(`
		html > 
		body > 
		table#table_useruimenu_10.tdContent >
		tbody >
		tr >
		td
	`)
	selection = selection.Next()
	selection = selection.Find(`
		form#frmUpload_4 >
		table#table_wu_1.tdContent >
		tbody >
		tr >
		td >
		table#table_wu_3 >
		tbody >
		tr >
		td >
		table#table_wu_19 >
		tbody >
		tr >
		input#trackid_1
	`)
	val, exists := selection.Attr("value")
	if !exists {
		return "", errors.New("Error: trackid not found")
	}

	fmt.Println("trackID: " + val)

	return val, nil
}
