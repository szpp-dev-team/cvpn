package util

import "regexp"

func CheckRegexp(pattern, s string) (bool, error) {
	return regexp.Match(pattern, []byte(s))
}
