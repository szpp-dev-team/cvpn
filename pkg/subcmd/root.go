package subcmd

import (
	"os"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cvpn",
		Short: "cvpn is a tool which makes you happy",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				_ = cmd.Help()
				os.Exit(0)
			}
		},
	}

	cmd.AddCommand(
		NewUploadCmd(),
		NewLoginCmd(),
		NewListCmd(),
		NewDownloadCmd(),
		NewFindCmd(),
		NewVersionCmd(),
	)

	return cmd
}
