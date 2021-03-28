package api

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
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
		"?trackid=" + trackID + strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10))
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

	{
		createField(mw, "xsauth", info.XSAuth)
		createField(mw, "txtServerUploadID", "")
	}

	// ファイルの書き出し
	{
		srcFile, err := os.Open(info.SrcPath)
		if err != nil {
			return "", nil, err
		}
		defer srcFile.Close()

		header := make(textproto.MIMEHeader)
		header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file1"; filename="%s"`, info.SrcPath))
		header.Set("Content-Type", contentTypeFromFile(srcFile))

		partWriter, err := mw.CreatePart(header)
		if err != nil {
			return "", nil, err
		}

		if _, err = io.Copy(partWriter, srcFile); err != nil {
			return "", nil, err
		}

		createField(mw, "txtRenameFile1", info.FileRenameOpt)
	}

	for i := 2; i <= 5; i++ {
		header := make(textproto.MIMEHeader)
		header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file%d"; filename=""`, i))
		header.Set("Content-Type", "application/octet-stream")
		_, _ = mw.CreatePart(header)

		createField(mw, "txtRenameFile"+strconv.Itoa(i), info.FileRenameOpt)
	}

	// フォームの普通のフィールドの書き出し
	{
		createField(mw, "t", "p")
		createField(mw, "v", info.VolumeID)
		createField(mw, "si", "")
		createField(mw, "ri", "")
		createField(mw, "pi", "")
		createField(mw, "dir", info.destDirPath)
		createField(mw, "acttype", "upload")
		createField(mw, "confirm", "yes")
		createField(mw, "trackid", info.trackID)
		createField(mw, "ignoreDfs", "1")
		createField(mw, "btnUpload", "アップロード")
	}

	if err := mw.Close(); err != nil {
		return "", nil, err
	}

	return mw.FormDataContentType(), formBodyBuff, nil
}

// https://qiita.com/ijufumi/items/c2d9f53262bb1f931d4e
func contentTypeFromFile(file *os.File) string {
	defer func() {
		_, _ = file.Seek(0, 0)
	}()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "application/octet-stream"
	}

	return http.DetectContentType(data)
}

func createField(mw *multipart.Writer, key, value string) {
	partWriter, _ := mw.CreateFormField(key)

	// queryedValue := url.QueryEscape(value)
	queryedValue := value
	_, _ = partWriter.Write([]byte(queryedValue))
}
