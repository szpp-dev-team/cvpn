package api

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type SegmentInfo struct {
	Name      string  `json:"name"`       // ファイル or ディレクトリ の名前
	Path      string  `json:"path"`       // ファイルならダウンロード URL、ディレクトリなら移動先の URL
	IsFile    bool    `json:"is_file"`    // file であるか
	IsDir     bool    `json:"is_dir"`     // dir であるか
	Size      float64 `json:"size"`       // サイズ
	Unit      string  `json:"unit"`       // サイズの単位
	UpdatedAt string  `json:"updated_at"` // できれば日時の構造体を使って欲しい
}

// セグメント情報の構造体のスライスを返す
func (c *Client) List(path, volumeID string) ([]SegmentInfo, error) {
	const ListEndpoint = "https://vpn.inf.shizuoka.ac.jp/dana/fb/smb/wfb.cgi"

	path = strings.Replace(path, "/", "\\", -1)

	params := genCommonAccessParam(volumeID, path)
	params.Set("sb", "name")
	params.Set("so", "asc")

	req, err := http.NewRequest(
		http.MethodGet,
		ListEndpoint,
		strings.NewReader(params.Encode()),
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.request(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // List() が終わる時に実行する
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("not 200")
	}

	segmentInfos, err := getSegmentInfos(resp.Body)
	if err != nil {
		return nil, err
	}

	return segmentInfos, nil
}

func getSegmentInfos(body io.ReadCloser) ([]SegmentInfo, error) {
	var segmentInfos []SegmentInfo

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	segmentInfos, err = findSegmentLines(doc)
	if err != nil {
		return nil, err
	}

	return segmentInfos, nil
}

func findSegmentLines(doc *goquery.Document) ([]SegmentInfo, error) {
	var segmentInfos []SegmentInfo

	selection := doc.Find("table#table_wfb_5 > tbody > script")

	lines := strings.Split(selection.Text()[1:], ";\n")

	for _, line := range lines {
		if len(line) == 0 {
			break
		}

		tokens := strings.Split(line[2:len(line)-1], ",")

		var tokensSeg SegmentInfo

		if len(tokens) == 3 { //ディレクトリの場合は要素数が3
			tokensSeg = SegmentInfo{
				Name:      tokens[0][1 : len(tokens[0])-1],
				Path:      tokens[1][1 : len(tokens[1])-1],
				IsDir:     true,
				Size:      -1,
				UpdatedAt: tokens[2][1 : len(tokens[2])-1],
			}
		}
		if len(tokens) == 4 { //ファイルの場合は要素数が4
			sizeItem := strings.Split(tokens[2][1:len(tokens[2])-1], "&")
			sizeValue, err := strconv.ParseFloat(sizeItem[0], 64)
			if err != nil {
				return nil, err
			}
			var sizeUnit string
			if sizeItem[1][len(sizeItem[1])-1] == 'B' { //最後がBとなっている場合はbytes以外
				sizeUnit = sizeItem[1][len(sizeItem[1])-2:]
			} else { //そうじゃない場合はbytes
				sizeUnit = sizeItem[1][len(sizeItem[1])-5:]
			}
			tokensSeg = SegmentInfo{
				Name:      tokens[0][1 : len(tokens[0])-1],
				Path:      tokens[1][1 : len(tokens[1])-1],
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
