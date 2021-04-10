package subcmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Shizuoka-Univ-dev/cvpn/static"
	"github.com/spf13/cobra"
)

func NewCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion",
		Short: "shell completion",
	}

	cmd.AddCommand(
		newBashCompletionCmd(),
	)

	return cmd
}

func newBashCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "bash",
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := BashCompletion()
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

func BashCompletion() (string, error) {
	indexFilePath, err := IndexFilePath()
	if err != nil {
		return "", err
	}

	completion := fmt.Sprintf(string(static.BashCompletionBytes), indexFilePath)

	return completion, nil
}
