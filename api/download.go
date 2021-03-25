// download api

package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// targetPath にあるファイルをダウンロードする。パスの指定は report を / としたもの。
func (c *Client) Download(targetPath, savePath, volumeID string) error {
	const DownloadURLFormat = "https://vpn.inf.shizuoka.ac.jp/dana/download/%s?url=/dana-cached/fb/smb/wfv.cgi?t=p&v=%s&si=0&ri=0&pi=0&ignoreDfs=1&dir=%s&file=%s"

	dirName, fileName := filepath.Split(targetPath)

	dir := strings.Replace(strings.Trim(dirName, "/"), "/", "\\", -1)

	params := genCommonAccessParam(volumeID, url.PathEscape(dir))
	params.Set("file", url.PathEscape(fileName))

	if savePath == "" {
		return errors.New("please specific save path")
	}

	return c.downloadFile(
		fmt.Sprintf(DownloadURLFormat, url.PathEscape(fileName), volumeID, url.PathEscape(dir), url.PathEscape(fileName)),
		fileName,
		savePath,
		params,
	)
}

// params を使いたかったんだけど、使うとなぜか permission denied でサーバーから弾かれてしまうので・・・。
func (c *Client) downloadFile(reqURL, fileName, dirPath string, params *url.Values) error {
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return err
	}

	resp, err := c.request(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("StatusCode of file downloading was %d (expected: 200 OK)", resp.StatusCode)
	}

	savePath := path.Join(dirPath, fileName)

	// TODO: 重複ケース
	file, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer file.Close()

	n, err := io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("Failed to download file(written 0 byte)")
	}

	return nil
}
