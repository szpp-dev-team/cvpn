package subcmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/Shizuoka-Univ-dev/cvpn/api"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func NewListCmd() *cobra.Command {
	var (
		volumeName  string
		showPathFlg bool
	)

	cmd := &cobra.Command{
		Use: "list",

		Aliases: []string{"ls"},
		Short:   "list files and directorys from path",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := godotenv.Load(); err != nil {
				return err
			}

			username := os.Getenv("SVPN_USERNAME")
			password := os.Getenv("SVPN_PASSWORD")

			client := api.NewClient()
			if err := client.LoadCookiesOrLogin(username, password); err != nil {
				return err
			}

			volumeID, err := api.GetVolumeIDFromName(volumeName)
			if err != nil {
				return err
			}

			segInfos, err := client.List(args[0], volumeID)
			if err != nil {
				return err
			}

			printSegmentInfos(segInfos, showPathFlg)

			return nil
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("write the path to place that you wnat to see all files and folders")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&volumeName, "volume", "v", api.VolumeNameFSShare, "volume id [fsshare / fs]")
	cmd.Flags().BoolVar(&showPathFlg, "path", false, "show segment's path")

	return cmd
}

func printSegmentInfos(segInfos []*api.SegmentInfo, showPathFlg bool) {
	for _, segInfo := range segInfos {
		if showPathFlg {
			if segInfo.IsFile {
				fmt.Printf("% 4.2f[%2s]   %s   %s(%s)\n", segInfo.Size, segInfo.Unit, segInfo.UpdatedAt, " "+segInfo.Name, segInfo.Path)
			} else {
				fmt.Printf("%8s   %s   %s(%s)\n", "-", segInfo.UpdatedAt, " "+segInfo.Name, segInfo.Path)
			}
		} else {
			if segInfo.IsFile {
				fmt.Printf("%-5.2f[%2s]   %s   %s\n", segInfo.Size, segInfo.Unit, segInfo.UpdatedAt, " "+segInfo.Name)
			} else {
				fmt.Printf("%9s   %s   %s\n", "-", segInfo.UpdatedAt, " "+segInfo.Name)
			}
		}
	}
}
