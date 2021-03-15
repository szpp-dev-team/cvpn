// list api

package api

import "net/http"

/*
TODO:
1. セグメント情報の構造体の定義
*/

// セグメント情報の構造体のスライスを返す
func (c *Client) List(path string) error {
	doc, err := c.getDoc(path, nil)
	if err != nil {
		return err
	}

	return nil
}