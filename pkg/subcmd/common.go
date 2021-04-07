package subcmd

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Shizuoka-Univ-dev/cvpn/api"
)

func Execute() {
	cmd := NewRootCmd()
	cmd.SetOutput(os.Stdout)

	if err := cmd.Execute(); err != nil {
		cmd.SetOutput(os.Stderr)
		cmd.Println(err)
		os.Exit(1)
	}

	if err := saveLogs(); err != nil {
		log.Fatal(err)
	}
}

func saveLogs() error {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	logDirPath := filepath.Join(userCacheDir, "cvpn", "log")
	_ = os.MkdirAll(logDirPath, 0755)

	logNameFormat := "2006-01-02_15-04-05.log"
	logName := time.Now().Format(logNameFormat)

	file, err := os.Create(filepath.Join(logDirPath, logName))
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(api.ReadLog()); err != nil {
		return err
	}

	return nil
}
