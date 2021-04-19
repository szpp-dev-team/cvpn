package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

//JSON用の構造体
type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//Username,passwordを返す（Config型）
func LoadConfig() (Config, error) {

	//ConfigDirPathを取得
	configFilePath, err := ConfigPath()
	if err != nil {
		return Config{}, err
	}

	jsonFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return Config{}, err
	}

	//Config型のデータ
	var data Config
	if err := json.Unmarshal(jsonFile, &data); err != nil {
		return Config{}, err
	}

	//username,password
	return data, nil
}

func ConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return path.Join(configDir, "cvpn", "config.json"), nil
}

func CreateConfigFile(username, password string) error {
	configFilePath, err := ConfigPath()
	if err != nil {
		return err
	}

	data := Config{
		Username: username,
		Password: password,
	}
	bytes, _ := json.Marshal(&data)

	if err := os.MkdirAll(filepath.Dir(configFilePath), 0700); err != nil {
		return err
	}
	return os.WriteFile(configFilePath, bytes, 0644)
}

func RemoveConfigFile() error {
	configPath, err := ConfigPath()
	if err != nil {
		return err
	}

	return os.RemoveAll(configPath)
}
