package static

import (
	_ "embed"
)

var (
	//go:embed completion/bash/cvpn
	BashCompletionBytes []byte
)
