package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		log.Println("WARNING: .env file not found, ", err)
	}

	err := godotenv.Load()
	if err != nil {
		log.Println("WARNING: Error loading .env file")
	}

	err = godotenv.Load(".env.test")
	if err != nil {
		log.Println("WARNING: No .env.test file found")
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
