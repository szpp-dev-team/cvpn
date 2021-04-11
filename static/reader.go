package static

import (
	_ "embed"
)

var (
	//go:embed completion/bash/cvpn
	BashCompletionBytes []byte

	//go:embed completion/zsh/cvpn
	ZshCompletionBytes []byte
)
