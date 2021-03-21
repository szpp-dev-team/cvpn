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
	var urls []SegmentInfo

	//要素をたどっていく
	selection := doc.Find(`
	html >
	body >
    table#table_useruimenu_10.tdContent >
    tbody >
	tr >
	td >
	form#theForm_3 >
    table#table_wfb_3 >
    tbody >
    tr >
    td >
    table#table_wfb_5 >
    tbody >
    script
	`)
	lines := strings.Split(selection.Text(), ";\n")

	first := true // なぜか最初の行だけ挙動がおかしいのでその修正

	for _, line := range lines {
		if len(line) == 0 {
			break
		}
		var tokens []string
		if first {
			tokens = strings.Split(line[3:len(line)-1], ",")
			first = false
		} else {
			tokens = strings.Split(line[2:len(line)-1], ",")
		}

		var tokens_seg SegmentInfo

		if len(tokens) == 3 { //ファイルの場合は要素数が3
			tokens_seg = SegmentInfo{
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
			Size_item := strings.Split(tokens[2][1:len(tokens[2])-1], "&")
			Size_value, _ := strconv.ParseFloat(Size_item[0], 64)
			var Size_unit string
			if Size_item[1][len(Size_item[1])-1] == 'B' { //最後がBとなっている場合はbytes以外
				Size_unit = Size_item[1][len(Size_item[1])-2:]
			} else {  //そうじゃない場合はbytes
				Size_unit = Size_item[1][len(Size_item[1])-5:]
			}
			tokens_seg = SegmentInfo{
				Name:      tokens[0][1 : len(tokens[0])-1],
				Path:      tokens[1][1 : len(tokens[1])-1],
				IsFile:    true,
				IsDir:     false,
				Size:      Size_value,
				Unit:      Size_unit,
				UpdatedAt: tokens[3][1 : len(tokens[3])-1],
			}
		}
		urls = append(urls, tokens_seg)
	}

	return urls, nil
}
