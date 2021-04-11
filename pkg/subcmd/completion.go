package subcmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Shizuoka-Univ-dev/cvpn/static"
	"github.com/spf13/cobra"
)

type shellEnum = int

const (
	// Start enum value from 1, NOT 0 (see: https://qiita.com/cia_rana/items/9d00ce81252ed970f362)
	bash shellEnum = iota + 1
	zsh
)

func NewCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion",
		Short: "shell completion",
	}

	cmd.AddCommand(
		newBashCompletionCmd(),
		newZshCompletionCmd(),
	)

	return cmd
}

func newBashCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "bash",
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := ShellCompletion(bash)
			if err != nil {
				return err
			}

			fmt.Println(string(b))

			return nil
		},
	}

	return cmd
}

func newZshCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "zsh",
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := ShellCompletion(zsh)
			if err != nil {
				return err
			}

			fmt.Println(string(b))

			return nil
		},
	}

	return cmd
}

func IndexFilePath() (string, error) {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	indexFilePath := filepath.Join(userCacheDir, "cvpn", "index.txt")
	return indexFilePath, nil
}

func ShellCompletion(shell shellEnum) (string, error) {
	indexFilePath, err := IndexFilePath()
	if err != nil {
		return "", err
	}

	var completion string
	switch shell {
	case bash:
		completion = fmt.Sprintf(string(static.BashCompletionBytes), indexFilePath)
	case zsh:
		completion = fmt.Sprintf(string(static.ZshCompletionBytes), indexFilePath)
	default:
		panic(shell)
	}

	return completion, nil
}
