package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Shizuoka-Univ-dev/cvpn/api"
	"github.com/Shizuoka-Univ-dev/cvpn/pkg/subcmd"
	"github.com/joho/godotenv"
)

func main() {
	subcmd.Execute()
}
