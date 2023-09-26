package subcmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var CvpnVersion string

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "show cvpn's version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("cvpn version %s %s/%s", CvpnVersion, runtime.GOOS, runtime.GOARCH)
		},
	}
}
