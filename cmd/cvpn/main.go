package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Shizuoka-Univ-dev/cvpn/api"
	"github.com/joho/godotenv"
)

func main() {

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
	if _, err := client.List(""); err != nil {
		log.Fatal(err)
	}

	fmt.Println("All success. exit 0.")
}
