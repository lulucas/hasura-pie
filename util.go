package pie

import (
	"math/rand"
	"os"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func StringPtr(s string) *string {
	return &s
}

func RestBaseUrl() string {
	baseUrl := strings.ToLower(os.Getenv("APP_REST_HOST"))
	if baseUrl == "" {
		baseUrl = "localhost"
		port := os.Getenv("APP_EXTERNAL_PORT")
		if port == "" {
			port = "8000"
		}
		if port != "80" {
			baseUrl += ":" + port
		}
	}
	if strings.ToLower(os.Getenv("APP_TLS_ENABLED")) == "true" {
		return "https://" + baseUrl
	} else {
		return "http://" + baseUrl
	}
}
