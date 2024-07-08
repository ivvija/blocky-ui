package settings

import (
	"os"

	"github.com/joho/godotenv"
)

var Host = "0.0.0.0"
var Port = "3000"
var ApiBaseUrl = "http://localhost:4000/api"

func init() {
	godotenv.Load()

	if host := os.Getenv("Host"); host != "" {
		Host = host
	}
	if port := os.Getenv("Port"); port != "" {
		Port = port
	}
	if url := os.Getenv("API_BASE_URL"); url != "" {
		ApiBaseUrl = url
	}
}
