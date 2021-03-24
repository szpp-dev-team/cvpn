package subcmd

import (
	"fmt"
	"log"
	"os"

	"github.com/Shizuoka-Univ-dev/cvpn/api"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "list",

		Aliases: []string{"ls"},
		Short:   "list files and directorys from path",
		Long:    "クソナガ説明",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				_ = cmd.Help()
				os.Exit(1)
			}
			// ls 処理
			if err := godotenv.Load(); err != nil {
				log.Fatal(err)
			}

			username := os.Getenv("SVPN_USERNAME")
			password := os.Getenv("SVPN_PASSWORD")

			fmt.Println(username, password)

			client := api.NewClient()
			if err := client.LoadCookiesOrLogin(username, password); err != nil {
				log.Fatal(err)
			}

			segInfos, err := client.List(args[0])
			if err != nil {
				log.Fatal(err)
			}

			s := ""

			fmt.Printf("index %4s |name %115s |size %4s |upload-at %4s\n", s, s, s, s)
			fmt.Println("ーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーーー")

			for i := 0; i < len(segInfos); i++ {
				fmt.Printf("%-11d|%-121s|%f\n", i+1, segInfos[i].Name, segInfos[i].Size)
			}
		},
	}

	return cmd
}
