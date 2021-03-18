package subcmd

import (
	"os"

	"github.com/spf13/cobra"
)

func NewLoginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "login to vpn service",
		Long:  "クソナガ説明",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.SetOut(os.Stderr)
			if len(args) > 0 {
				cmd.Println("too many args")
				_ = cmd.Help()
				os.Exit(1)
			}

			// ログイン処理
		},
	}
}
