package utils

import (
	"fmt"
	"os"
	"strings"
)

func GetFrontendURL() string {
	envURLs := os.Getenv("FRONTEND_URLS")
	// split by comma
	frontendEnvURLs := strings.Split(envURLs, ",")

	mode := os.Getenv("APP_MODE")
	fmt.Println("frontendEnvURLs", frontendEnvURLs)

	if mode == "production" {
		return frontendEnvURLs[1]
	} else {
		return frontendEnvURLs[0]
	}
}
