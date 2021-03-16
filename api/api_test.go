package api

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestDownload(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatal(err)
	}

	client := NewClient()
	if err := client.LoadCookiesOrLogin(os.Getenv("SVPN_USERNAME"), os.Getenv("SVPN_PASSWORD")); err != nil {
		t.Fatal(err)
	}

	const (
		targetPath = "/path/to/file"
		savePath   = ""
		volumeID   = VolumeIDFSShare
	)
	if err := client.Download(
		targetPath,
		savePath,
		volumeID,
	); err != nil {
		t.Fatal(err)
	}
}

func TestUpload(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatal(err)
	}

	client := NewClient()
	if err := client.LoadCookiesOrLogin(os.Getenv("SVPN_USERNAME"), os.Getenv("SVPN_PASSWORD")); err != nil {
		t.Fatal(err)
	}

	const (
		srcFilePath     = "./../test.txt"
		renamedFileName = "upload-rename-test.txt"
		volumeID        = "resource_1389773645.177066.2,2020"
		destDirPath     = "cs200XX"
	)
	if err := client.UploadFile(srcFilePath, renamedFileName, volumeID, destDirPath); err != nil {
		log.Fatal(err)
	}
}
