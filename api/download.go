// download api

package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// targetPath にあるファイルをダウンロードする。パスの指定は report を / としたもの。
// savePath が {dir} だったらサーバー上と同名でそこに保存し、{dir}/{file} だったら {dir} 上に名前は {file} で保存する。
// savePath の指定がなければカレントディレクトリ上に保存する。
func (c *Client) Download(targetPath, savePath, volumeID string) error {
	const DownloadURLFormat = "https://vpn.inf.shizuoka.ac.jp/dana/download/%s?url=/dana-cached/fb/smb/wfv.cgi?v=%s&dir=%s&file=%s"

	dirName, fileName := filepath.Split(targetPath)

	dir := strings.Replace(strings.Trim(dirName, "/"), "/", "\\", -1)

	params := genCommonAccessParam(volumeID, url.PathEscape(dir))
	params.Set("file", url.PathEscape(fileName))

	if savePath == "" {
		tmp, err := os.Getwd()
		if err != nil {
			return err
		}
		savePath = tmp + "/" + fileName
	}

	return c.downloadFile(
		fmt.Sprintf(DownloadURLFormat, url.PathEscape(fileName), volumeID, url.PathEscape(dir), url.PathEscape(fileName)),
		savePath,
		params,
	)
}

// params を使いたかったんだけど、使うとなぜか permission denied でサーバーから弾かれてしまうので・・・。
func (c *Client) downloadFile(reqURL, savePath string, params *url.Values) error {
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
