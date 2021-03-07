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
	if err := client.Login(username, password); err != nil {
		log.Fatal(err)
	}
}
