package settings

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var Host = "0.0.0.0"
var Port = "3000"
var ApiBaseUrl = "http://localhost:4000/api"
var PauseDuration = time.Minute * 5

func init() {
	_ = godotenv.Load()

	if host := os.Getenv("HOST"); host != "" {
		Host = host
	}
	if port := os.Getenv("PORT"); port != "" {
		Port = port
	}
	if url := os.Getenv("API_BASE_URL"); url != "" {
		ApiBaseUrl = url
	}
	if duration := os.Getenv("PAUSE_DURATION"); duration != "" {
		d, err := time.ParseDuration(duration)
		if err != nil {
			log.Println("Error parsing PAUSE_DURATION, using default:", err)
		} else {
			PauseDuration = d
		}
	}
}
