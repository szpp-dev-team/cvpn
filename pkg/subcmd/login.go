package subcmd

import (
	"log"
	"os"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/Shizuoka-Univ-dev/cvpn/api"
	"github.com/Shizuoka-Univ-dev/cvpn/pkg/util"
	"encoding/json"
	"path"
	
)

//JSON用の構造体
type Config struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
}

func NewLoginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "login to vpn service",
		Long:  "クソナガ説明",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.SetOut(os.Stderr)
			if len(args) > 0 {
				cmd.Println("too many args")
				_ = cmd.Help()
				os.Exit(1)
			}

			// ログイン処理
			var username,password string

			//ConfigDirPathを取得
			configDir , err := os.UserConfigDir()
			if err != nil {
				log.Fatal(err)
			}

			// ConfigFilePath（configFileを書き込むパス）の設定
			configFilePath := path.Join(configDir,"cvpn/config.json") 

			//入力
			fmt.Print("username >> ")
			fmt.Scan(&username)

			fmt.Print("password >> ")
			fmt.Scan(&password)

			//接続
			client := api.NewClient()

			// ログイン処理
			if err := client.Login(username,password); err != nil {
				fmt.Println("Either the username or password is invalid.")
				log.Fatal(err)
			}else{
				
				// 生成確認
				if flag,err := util.InputYN("Creating configFile? [Y/n] >> "); flag && err==nil {

					if err := os.MkdirAll(path.Dir(configFilePath), 0700); err != nil {
						log.Fatal(err)
					}

					fp, err := os.OpenFile(configFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
					if err != nil {
						log.Fatal(err)
					}

					defer fp.Close()
					
					//JSONデータ
					data := Config{
						Username:     username,
						Password:     password,
					}
					
					bytes, _ := json.Marshal(&data)
					
					_,err = fp.WriteString(string(bytes))
					if err != nil {
						log.Fatal(err)
					}

					// ファイル生成（更新）ログ
					log.Printf("Created configFile into %q.\n",configFilePath)
					
				}else{

					// ファイル生成（更新）中止ログ
					log.Printf("Not created configFile.\n")
				}
					
			}

		},
	}
}

   