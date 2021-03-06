package api

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type SegmentInfo struct {
	Name      string  `json:"name"`    // ファイル or ディレクトリ の名前
	Path      string  `json:"path"`    // ファイルならダウンロード URL、ディレクトリなら移動先の URL
	IsFile    bool    `json:"is_file"` // file であるか
	IsDir     bool    `json:"is_dir"`  // dir であるか
	Size      float64 `json:"size"`    // サイズ
	Unit      string  `json:"unit"`    // サイズの単位
	VolumeID  string  `json:"volume_id"`
	UpdatedAt string  `json:"updated_at"` // TODO: できれば日時の構造体を使って欲しい
}

type PermissionError struct {
	RequestPath string
	VolumeID    string
}

func (pe *PermissionError) Error() string {
	return "Perrmission Denied"
}

// セグメント情報の構造体のスライスを返す
func (c *Client) List(path, volumeID string) ([]*SegmentInfo, error) {
	const ListEndpoint = "https://vpn.inf.shizuoka.ac.jp/dana/fb/smb/wfb.cgi"

	escapedPath := url.QueryEscape(path)
	escapedPath = strings.Replace(escapedPath, "/", "\\", -1)

	params := fmt.Sprintf(
		"?t=p&v=%s&si=0&ri=0&pi=0&sb=%s&so=%s&dir=%s",
		volumeID,
		"name",
		"asc",
		escapedPath,
	)

	req, err := http.NewRequest(
		http.MethodGet,
		ListEndpoint+params,
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.request(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusFound { // permission denied
		return nil, &PermissionError{RequestPath: path, VolumeID: volumeID}
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("StatusCode of file uploading was %d (expected: 200 OK)", resp.StatusCode)
	}

	segmentInfos, err := getSegmentInfos(b, path)
	if err != nil {
		return nil, err
	}

	return segmentInfos, nil
}

func IsPermissionDenied(err error) bool {
	switch err.(type) {
	case *PermissionError:
		return true
	default:
		return false
	}
}

func getSegmentInfos(b []byte, dirPath string) ([]*SegmentInfo, error) {
	var segmentInfos []*SegmentInfo

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	segmentInfos, err = findSegmentLines(doc)
	if err != nil {
		return nil, err
	}

	for _, segmentInfo := range segmentInfos {
		s := path.Join(dirPath, segmentInfo.Name)
		segmentInfo.Path = s
	}

	return segmentInfos, nil
}

func findSegmentLines(doc *goquery.Document) ([]*SegmentInfo, error) {
	var segmentInfos []*SegmentInfo

	selection := doc.Find("table#table_wfb_5 > tbody > script")
	if selection.Length() == 0 {
		return nil, errors.New("the directory name is wrong, or the directory have no data")
	}

	lines := strings.Split(selection.Text()[1:], ";\n")

	for _, line := range lines {
		if len(line) == 0 {
			break
		}

		tokens, err := splitSegmentLine(line[2 : len(line)-1])
		if err != nil {
			return nil, err
		}

		var tokensSeg *SegmentInfo
		if len(tokens) == 3 { //ディレクトリの場合は要素数が3
			tokensSeg = &SegmentInfo{
				Name:      tokens[0][1 : len(tokens[0])-1],
				IsDir:     true,
				Size:      -1,
				UpdatedAt: tokens[2][1 : len(tokens[2])-1],
			}
		} else if len(tokens) == 4 { //ファイルの場合は要素数が4
			sizeItem := strings.Split(tokens[2][1:len(tokens[2])-1], "&")
			sizeValue, sizeUnit := fileSize(sizeItem)

			tokensSeg = &SegmentInfo{
				Name:      tokens[0][1 : len(tokens[0])-1],
				IsFile:    true,
				Size:      sizeValue,
				Unit:      sizeUnit,
				UpdatedAt: tokens[3][1 : len(tokens[3])-1],
			}
		}

		segmentInfos = append(segmentInfos, tokensSeg)
	}

	return segmentInfos, nil
}

func fileSize(sizeItem []string) (float64, string) {
	sizeValue, err := strconv.ParseFloat(sizeItem[0], 64)
	if err != nil {
		return 0, ""
	}

	var sizeUnit string
	if sizeItem[1][len(sizeItem[1])-1] == 'B' { //最後がBとなっている場合はbytes以外
		sizeUnit = sizeItem[1][len(sizeItem[1])-2:]
	} else { //そうじゃない場合はbytes
		sizeUnit = "B"
	}

	return sizeValue, sizeUnit
}

func splitSegmentLine(line string) ([]string, error) {
	r, err := regexp.Compile(`[^\s"']+|"([^"]*)"|'([^']*)`)
	if err != nil {
		return nil, err
	}

	tokens := r.FindAllString(line, -1)
	var ans []string
	for _, token := range tokens {
		if token != "," {
			ans = append(ans, token)
		}
	}

	return ans, nil
}
