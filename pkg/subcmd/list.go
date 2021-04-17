package subcmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/Shizuoka-Univ-dev/cvpn/api"
	"github.com/Shizuoka-Univ-dev/cvpn/pkg/config"
	"github.com/spf13/cobra"
)

func NewListCmd() *cobra.Command {
	var (
		volumeName  string
		showPathFlg bool
		jsonFlg     bool
	)

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list files and directorys from path",
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := config.LoadConfig()
			if err != nil {
				return err
			}
			client := api.NewClient()
			if err := client.LoadCookiesOrLogin(config.Username, config.Password); err != nil {
				return err
			}
			volumeID, err := api.GetVolumeIDFromName(volumeName)
			if err != nil {
				return err
			}

			segInfos, err := client.List(args[0], volumeID)
			if err != nil {
				if api.IsPermissionDenied(err) {
					fmt.Fprintf(os.Stderr, "cannot fetch directory '%s': Permission denied", args[0])
					return nil
				}
				return err
			}

			printSegmentInfos(segInfos, showPathFlg, jsonFlg)

			return nil
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("write the path to place that you want to see all files and folders")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&volumeName, "volume", "v", api.VolumeNameFSShare, "volume id [fsshare / fs]")
	cmd.Flags().BoolVar(&showPathFlg, "path", false, "show path of segment")
	cmd.Flags().BoolVar(&jsonFlg, "json", false, "show segments in json format")

	return cmd
}

func printSegmentInfos(segInfos []*api.SegmentInfo, showPathFlg, jsonFlg bool) {
	if jsonFlg {
		json, _ := json.Marshal(segInfos)
		fmt.Println(string(json))
		return
	}

	for _, segInfo := range segInfos {
		if showPathFlg {
			if segInfo.IsFile {
				fmt.Printf("%-6.2f[%2s]   %s   %s(%s)\n", segInfo.Size, segInfo.Unit, segInfo.UpdatedAt, " "+segInfo.Name, segInfo.Path)
			} else {
				fmt.Printf("%10s   %s   %s(%s)\n", "-", segInfo.UpdatedAt, " "+segInfo.Name, segInfo.Path)
			}
		} else {
			if segInfo.IsFile {
				fmt.Printf("%-6.2f[%2s]   %s   %s\n", segInfo.Size, segInfo.Unit, segInfo.UpdatedAt, " "+segInfo.Name)
			} else {
				fmt.Printf("%10s   %s   %s\n", "-", segInfo.UpdatedAt, " "+segInfo.Name)
			}
		}
	}
}
