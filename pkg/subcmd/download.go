package subcmd

import (
	"errors"
	"strings"

	"github.com/Shizuoka-Univ-dev/cvpn/api"
	"github.com/Shizuoka-Univ-dev/cvpn/pkg/config"
	"github.com/spf13/cobra"
)

func NewDownloadCmd() *cobra.Command {
	var (
		savePath   string
		volumeName string
	)

	cmd := &cobra.Command{
		Use:   "download",
		Short: "download a file from path",
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

			pathes := strings.Split(args[0], "\n")
			for _, path := range pathes {
				if err := client.Download(path, savePath, volumeID); err != nil {
					return err
				}
			}

			return nil
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires target path")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&savePath, "output", "o", ".", "a path for saving downloaded files")
	cmd.Flags().StringVarP(&volumeName, "volume", "v", api.VolumeNameFSShare, "volume id [fsshare / fs]")

	return cmd
}
