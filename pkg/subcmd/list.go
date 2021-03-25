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
	var volumeName string

	cmd := &cobra.Command{
		Use: "list",

		Aliases: []string{"ls"},
		Short:   "list files and directorys from path",
		Long:    "クソナガ説明",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := godotenv.Load(); err != nil {
				return err
			}

			username := os.Getenv("SVPN_USERNAME")
			password := os.Getenv("SVPN_PASSWORD")

			fmt.Println(username, password)

			client := api.NewClient()
			if err := client.LoadCookiesOrLogin(username, password); err != nil {
				return err
			}

			// TODO
			_, err := api.GetVolumeIDFromName(volumeName)
			if err != nil {
				return err
			}

			segInfos, err := client.List(args[0])
			if err != nil {
				return err
			}

			s := ""

			//format := genFormat(segInfos)　保留

			fmt.Printf("index %4s name %50s size %4s upload-at %4s\n", s, s, s, s)
			fmt.Println("ーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーー")

			for i, _ := range segInfos {
				if segInfos[i].Size == -1 {
					fmt.Printf("%-11d%-50s  %s %s%s\n", i+1, segInfos[i].Name, segInfos[i].Path, s, segInfos[i].UpdatedAt)
				} else {
					f := fmt.Sprintf("%.2f", segInfos[i].Size)
					fmt.Printf("%-11d%-50s  %s %s%s %s\n", i+1, segInfos[i].Name, segInfos[i].Path, f, segInfos[i].Unit, segInfos[i].UpdatedAt)
				}
				/*fmt.Println(segInfos[i].Unit)
				length := countStrLen(segInfos[i].Name)
				if segInfos[i].Size == -1 {
					fmt.Printf(format, "\n", i+1, segInfos[i].Name, segInfos[i].Path, s, segInfos[i].UpdatedAt)
				} else {
					f := fmt.Sprintf("%.2f", segInfos[i].Size)
					fmt.Printf(format, "\n", i+1, segInfos[i].Name, f, segInfos[i].Unit, segInfos[i].UpdatedAt)
				} 保留 */
			}

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

	return cmd
}

/*func genFormat(s []api.SegmentInfo) string {
	format := `%%-%dd%%-%ds%%-%ds %%s %%s`
	maxname := 0
	maxpath := 0
	for i, _ := range s {
		if x1 := countStrLen(s[i].Name); maxname < x1 {
			maxname = x1
		}
		if x2 := countStrLen(s[i].Path); maxpath < x2 {
			maxpath = x2
		}
	}
	a := fmt.Sprintf(format, 11, maxname, maxpath)
	fmt.Println(a)
	return a
}

func countStrLen(s string) int {
	length := 0

	for _, c := range s {
		if '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' {
			length += 1
		} else {
			length += 2
		}
	}
	return length
}　　保留*/
