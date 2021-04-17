package subcmd

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/Shizuoka-Univ-dev/cvpn/api"
	"github.com/Shizuoka-Univ-dev/cvpn/pkg/config"
	"github.com/spf13/cobra"
)

var (
	volumeName   string
	namePattern  string
	pathPattern  string
	recursiveFlg bool
)

func NewFindCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find",
		Short: "find files or directories which matches pattern",
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

			_, err = searchMatches(
				client,
				args[0],
				volumeID,
				func(s string) {
					fmt.Println(s)
				},
			)
			if err != nil {
				return err
			}

			return nil
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("{starting-directory} was not set")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&volumeName, "volume", "v", api.VolumeNameFSShare, "volume id [fsshare / fs]")
	cmd.Flags().BoolVarP(&recursiveFlg, "recursive", "r", false, "recursively search files or directories")
	cmd.Flags().StringVar(&namePattern, "name", "", "pattern for name of files or directories")
	cmd.Flags().StringVar(&pathPattern, "path", "", "pattern for path of files or directories")

	return cmd
}

func checkRegexp(pattern, s string) (bool, error) {
	return regexp.Match(pattern, []byte(s))
}

func searchMatches(client *api.Client, begin, volumeID string, matchHandler func(string)) ([]string, error) {
	var stack, matchedPathes []string
	stack = append(stack, begin)
	for len(stack) > 0 {
		tail := len(stack) - 1
		segs, err := client.List(stack[tail], volumeID)
		if err != nil {
			if api.IsPermissionDenied(err) {
				fmt.Fprintf(os.Stderr, "cannot fetch directory '%s': Permission denied\n", stack[tail])
				stack = stack[:tail]
				continue
			} else {
				return nil, err
			}
		}
		stack = stack[:tail]

		for _, seg := range segs {
			ok1, err := checkRegexp(namePattern, seg.Name)
			if err != nil {
				return nil, err
			}
			ok2, err := checkRegexp(pathPattern, seg.Path)
			if err != nil {
				return nil, err
			}
			if ok1 && ok2 {
				matchHandler(seg.Path)
				matchedPathes = append(matchedPathes, seg.Path)
			}
			if seg.IsDir && recursiveFlg {
				stack = append(stack, seg.Path)
			}
		}
	}

	return matchedPathes, nil
}
