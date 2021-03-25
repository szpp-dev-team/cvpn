package subcmd

import (
	"os"

	"github.com/spf13/cobra"
)

func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "list",

		Aliases: []string{"ls"},
		Short:   "list files and directorys from path",
		Long:    "クソナガ説明",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				_ = cmd.Help()
				os.Exit(1)
			}

			// ls 処理
		},
	}

	return cmd
}
