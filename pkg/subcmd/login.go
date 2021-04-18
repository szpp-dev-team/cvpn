package subcmd

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Shizuoka-Univ-dev/cvpn/api"
	"github.com/Shizuoka-Univ-dev/cvpn/pkg/config"
	"github.com/Shizuoka-Univ-dev/cvpn/pkg/util"
	"github.com/spf13/cobra"
)

func NewLoginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "login to vpn service",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.SetOut(os.Stderr)

			client := api.NewClient()
			if ok := client.CheckCookies(); ok {
				yes, err := util.InputYN("You seem to have already logined. Do you want to login again?[Y/n] ")
				if err != nil {
					log.Fatal(err)
				}
				if !yes {
					fmt.Println("login canceled")
					return
				}
				if err := client.Logout(); err != nil {
					log.Fatal(err)
				}
			}

			if err := config.RemoveConfigFile(); err != nil {
				log.Fatal(err)
			}
			
			username, password, err := InputUserInfo(client)
			if err != nil {
				fmt.Println("Either the username or password is invalid.")
				log.Fatal(err)
			}

			flag, err := util.InputYN("Creating configFile? [Y/n] >> ")
			if err != nil {
				log.Fatal(err)
			}
			if flag {
				if err := config.CreateConfigFile(username, password); err != nil {
					log.Fatal(err)
				}
				log.Println("Created config file.")
			} else {
				log.Printf("Not created configFile.\n")
			}
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return errors.New("too many args")
			}

			return nil
		},
	}
}

// id と ps を入力して認証に失敗したら error
func InputUserInfo(client *api.Client) (string, string, error) {
	var username, password string

	fmt.Print("username >> ")
	fmt.Scan(&username)
	fmt.Print("password >> ")
	fmt.Scan(&password)

	if err := client.Login(username, password); err != nil {
		return "", "", err
	}

	return username, password, nil
}
