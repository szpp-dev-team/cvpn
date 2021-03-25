package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

//JSON用の構造体
type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//Username,passwordを返す（Config型）
func LoadConfig() (Config, error) {

	//ConfigDirPathを取得
	configDir, err := os.UserConfigDir()
	if err != nil {
		return Config{}, err
	}

	// ConfigFilePath（configFileを書き込むパス）の設定
	configFilePath := path.Join(configDir, "cvpn/config.json")

	jsonFile, err := ioutil.ReadFile(configFilePath)

	if err != nil {
		return Config{}, err
	}

	//Config型のデータ
	var data Config
	err = json.Unmarshal(jsonFile, &data)
	if err != nil {
		return Config{}, err
	}

	//username,password
	return data, nil
}
