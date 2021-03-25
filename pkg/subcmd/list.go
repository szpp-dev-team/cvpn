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

			printSegmentInfo(segInfos, showPathFlg)

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

func printSegmentInfo(segInfos []*api.SegmentInfo, showPathFlg bool) {
	if showPathFlg {
		fmt.Printf("%-7s%-50s%50s%12s%8s\n", "index", "name", "path", "size", "uploaded_at")
	} else {
		fmt.Printf("%-7s%-50s%-14s%-8s\n", "index", "name", "size", "uploaded_at")
	}

	for i, segInfo := range segInfos {
		if showPathFlg {
			if segInfo.IsDir {
				fmt.Printf("%-6d %-50s %s %-12s %-8s\n", i+1, " "+segInfo.Name, segInfo.Path, "", segInfo.UpdatedAt)
			} else {
				fmt.Printf("%-6d %-50s %s %.2f[%s] %-8s\n", i+1, " "+segInfo.Name, segInfo.Path, segInfo.Size, segInfo.Unit, segInfo.UpdatedAt)
			}
		} else {
			if segInfo.IsDir {
				fmt.Printf("%-6d %-50s %-12s %-8s\n", i+1, " "+segInfo.Name, "", segInfo.UpdatedAt)
			} else {
				fmt.Printf("%-6d %-50s %.2f[%s] %-8s\n", i+1, " "+segInfo.Name, segInfo.Size, segInfo.Unit, segInfo.UpdatedAt)
			}
		}
	}
}
