package initializers

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

var EAT *time.Location

func InitTimezone() {
	var err error
	EAT, err = time.LoadLocation("Africa/Nairobi")
	if err != nil {
		log.Fatalf("failed to load timezone: %v", err)
	}
}
