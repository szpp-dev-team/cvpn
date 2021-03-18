package api

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestLogin(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatal(err)
	}

	client := NewClient()
	if err := client.LoadCookiesOrLogin(os.Getenv("SVPN_USERNAME"), os.Getenv("SVPN_PASSWORD")); err != nil {
		switch err.(type) {
		case *SessionError:
			// セッションエラーだからセッションを選ばせればいい
		default:
			t.Fatal(err)
		}
	}
}

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
		log.Fatal(err)
	}
}
