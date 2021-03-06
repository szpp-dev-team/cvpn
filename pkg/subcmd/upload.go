package subcmd

import (
	"errors"
	"log"
	"os"

	"github.com/Shizuoka-Univ-dev/cvpn/api"
	"github.com/Shizuoka-Univ-dev/cvpn/pkg/config"
	"github.com/spf13/cobra"
)

func NewUploadCmd() *cobra.Command {

	var volumeName string

	cmd := &cobra.Command{
		Use:   "upload",
		Short: "upload a file from path",
		Run: func(cmd *cobra.Command, args []string) {
			srcPath := args[0]    //ソース
			uploadPath := args[1] //アップロード先のURL

			client := api.NewClient()
			config, err := config.LoadConfig()
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			if err := client.LoadCookiesOrLogin(config.Username, config.Password); err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			var rename string = "" // TODO
			volumeID, err := api.GetVolumeIDFromName(volumeName)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			if err := client.UploadFile(srcPath, rename, volumeID, uploadPath); err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("Not enough args")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&volumeName, "volume", "v", api.VolumeIDFSShare, "volume id [fsshare / fs]")

	return cmd
}
