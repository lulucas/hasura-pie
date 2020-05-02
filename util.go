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
	if strings.ToLower(os.Getenv("APP_TLS_ENABLED")) == "true" {
		return "https://" + strings.ToLower(os.Getenv("APP_REST_HOST"))
	} else {
		return "http://" + strings.ToLower(os.Getenv("APP_REST_HOST"))
	}
}
