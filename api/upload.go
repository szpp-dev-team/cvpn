package api

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type uploadReqInfo struct {
	SrcPath       string
	FileRenameOpt string // アップロード後に名前を変えるかどうか。変えないなら空文字列。
	VolumeID      string
	destDirPath   string
	trackID       string
	XSAuth        string
}

// ローカルの `srcPath` を `volumeID` ドライブの `destDirPath` へアップロードする。
// `fileRenameOpt` が空文字列でなければ、アップロードしたファイルの名前は `fileRenameOpt` へ変更される。
// `destDirPath` の '/' は '\' に置換されて送信される。
func (c *Client) UploadFile(srcPath, fileRenameOpt, volumeID, destDirPath string) error {
	destDirPath = strings.Replace(destDirPath, "/", "\\", -1)

	trackID, err := c.scrapeTrackID4Upload(volumeID, destDirPath)
	if err != nil {
		return err
	}

	contentType, multipartBody, err := createMultipartBody4Upload(&uploadReqInfo{
		SrcPath:       srcPath,
		FileRenameOpt: fileRenameOpt,
		VolumeID:      volumeID,
		destDirPath:   destDirPath,
		trackID:       trackID,
		XSAuth:        c.authParams.Get("xsauth"),
	})
	if err != nil {
		return err
	}

	// アップロード先URL には trackID の末尾にUNIX時間を追加したクエリパラメータを付与する
	uploadURL := ("https://vpn.inf.shizuoka.ac.jp/dana/fb/smb/wu.cgi" +
		"?trackid=" + trackID + strconv.FormatInt(time.Now().Unix(), 10))
	req, err := http.NewRequest(http.MethodPost, uploadURL, multipartBody)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := c.request(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("StatusCode of file uploading was %d (expected: 200 OK)", resp.StatusCode)
	}
	return nil
}

func (c *Client) scrapeTrackID4Upload(volumeID, destDirPath string) (string, error) {
	const uploadFormURL = "https://vpn.inf.shizuoka.ac.jp/dana/fb/smb/wu.cgi"

	xsauth := c.authParams.Get("xsauth")
	if len(xsauth) == 0 {
		return "", fmt.Errorf("Client has no `xsauth`")
	}
	postForm := genCommonAccessParam(volumeID, destDirPath)
	postForm.Set("acttype", "upload")
	postForm.Set("ignoreDfs", "1")
	postForm.Set("xsauth", xsauth)

	req, err := http.NewRequest(http.MethodPost, uploadFormURL, strings.NewReader(postForm.Encode()))
	if err != nil {
		return "", err
	}
	resp, err := c.request(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	docScraper, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}
	trackid, exists := docScraper.Find("#trackid_1").Attr("value")
	if !exists {
		return "", fmt.Errorf("`trackid` not found: request URL: %q", uploadFormURL)
	}
	return trackid, nil
}

// https://medium.com/eureka-engineering/multipart-file-upload-in-golang-c4a8eb15a3ee
func createMultipartBody4Upload(info *uploadReqInfo) (contentType string, body *bytes.Buffer, err error) {
	formBodyBuff := &bytes.Buffer{}
	mw := multipart.NewWriter(formBodyBuff)

	// ファイルの書き出し
	{
		srcFile, err := os.Open(info.SrcPath)
		if err != nil {
			return "", nil, err
		}
		defer srcFile.Close()

		partWriter, _ := mw.CreateFormFile("file1", path.Base(info.SrcPath))
		if _, err = io.Copy(partWriter, srcFile); err != nil {
			return "", nil, err
		}
	}

	// フォームの普通のフィールドの書き出し
	{
		formFields := genCommonAccessParam(info.VolumeID, info.destDirPath)
		formFields.Set("acttype", "upload")
		formFields.Set("confirm", "yes")
		formFields.Set("ignoreDfs", "1")
		formFields.Set("trackid", info.trackID)
		formFields.Set("txtRenameFile1", info.FileRenameOpt)
		formFields.Set("txtServerUploadID", "")
		formFields.Set("xsauth", info.XSAuth)

		for key, values := range *formFields {
			partWriter, err := mw.CreateFormField(key)
			if err != nil {
				return "", nil, err
			}
			if len(values) > 1 {
				return "", nil, fmt.Errorf("The form value of %q cannot be multiple (got: %v)", key, values)
			}
			value := url.QueryEscape(formFields.Get(key))
			if _, err = partWriter.Write([]byte(value)); err != nil {
				return "", nil, err
			}
		}
	}
	if err := mw.Close(); err != nil {
		return "", nil, err
	}
	return mw.FormDataContentType(), formBodyBuff, nil
}
