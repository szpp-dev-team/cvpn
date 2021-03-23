// list api

package api

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//フィールドも先頭一文字が大文字かどうかで public かどうかが決まる
type SegmentInfo struct {
	Name      string // ファイル or ディレクトリ の名前
	Path      string // ファイルならダウンロード URL、ディレクトリなら移動先の URL
	IsFile    bool   // 良いデザインパターンがあるはずなのであったらそれを採用してください
	IsDir     bool
	Size      float64 // サイズ
	Unit      string  // サイズの単位
	UpdatedAt string  // できれば日時の構造体を使って欲しい
}

// セグメント情報の構造体のスライスを返す
func (c *Client) List(path string) ([]SegmentInfo, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://vpn.inf.shizuoka.ac.jp/dana/fb/smb/wfb.cgi?t=p&v=resource_1423533946.487706.3&si=0&ri=0&pi=0&sb=name&so=asc&dir=report",
		nil,
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

	// resp.Body は html
	segmentInfos, err := getSegmentInfos(resp.Body)
	if err != nil {
		return nil, err
	}

	return segmentInfos, nil
}

// func関数名(引数) (戻り値) {}
//戻り値は2つ以上なら () をつける
func getSegmentInfos(body io.ReadCloser) ([]SegmentInfo, error) {
	var segmentInfos []SegmentInfo

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	_, err = findSegmentLines(doc)
	if err != nil {
		return nil, err
	}

	// ここで lines からファイル名とかサイズとかを抜きとり、SegmentInfo のスライスで返す

	return segmentInfos, nil
}

// "d(...)" とか "f(...)"とかの形式で返す
func findSegmentLines(doc *goquery.Document) ([]SegmentInfo, error) {
	var segmentInfos []SegmentInfo

	//要素をたどっていく
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
				IsFile:    false,
				IsDir:     true,
				Size:      -1,
				Unit:      "",
				UpdatedAt: tokens[2][1 : len(tokens[2])-1],
			}
		}
		if len(tokens) == 4 { //ファイルの場合は要素数が4
			sizeItem := strings.Split(tokens[2][1:len(tokens[2])-1], "&")
			sizeValue, err := strconv.ParseFloat(sizeItem[0], 64)
			if err != nil {
				return nil, err
			}
			var Size_unit string
			if sizeItem[1][len(sizeItem[1])-1] == 'B' { //最後がBとなっている場合はbytes以外
				Size_unit = sizeItem[1][len(sizeItem[1])-2:]
			} else { //そうじゃない場合はbytes
				Size_unit = sizeItem[1][len(sizeItem[1])-5:]
			}
			tokensSeg = SegmentInfo{
				Name:      tokens[0][1 : len(tokens[0])-1],
				Path:      tokens[1][1 : len(tokens[1])-1],
				IsFile:    true,
				IsDir:     false,
				Size:      sizeValue,
				Unit:      Size_unit,
				UpdatedAt: tokens[3][1 : len(tokens[3])-1],
			}
		}

		segmentInfos = append(segmentInfos, tokensSeg)
	}

	return segmentInfos, nil
}
