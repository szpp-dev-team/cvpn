package api

import (
	"fmt"
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
		targetPath = "/cs20097/makabe.png"
		savePath   = ""
		volumeID   = VolumeIDFS + "2020"
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
		volumeID        = VolumeIDFS + ",{dir}"
		destDirPath     = "cs200XX"
	)
	if err := client.UploadFile(srcFilePath, renamedFileName, volumeID, destDirPath); err != nil {
		t.Fatal(err)
	}
}

func TestList(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatal(err)
	}

	client := NewClient()
	if err := client.LoadCookiesOrLogin(os.Getenv("SVPN_USERNAME"), os.Getenv("SVPN_PASSWORD")); err != nil {
		t.Fatal(err)
	}

	const (
		path = "/report"
		// volumeID = VolumeIDFSShare
	)
	segmentInfos, err := client.List(path)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(segmentInfos)
}
