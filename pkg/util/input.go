package util

import (
	"fmt"
	"strings"
)

// prompt を表示して y か n を入力させる。y なら true, n なら false
// 入力が1文字でなければ再入力させる。
func InputYN(prompt string) (bool, error) {
	for {
		fmt.Print(prompt + " ")

		var yn string
		if _, err := fmt.Scan(&yn); err != nil {
			return false, err
		}

		if len(yn) != 1 {
			continue
		}

		if strings.ToLower(yn) == "y" {
			return true, nil
		} else if strings.ToLower(yn) == "n" {
			return false, nil
		}
	}
}
