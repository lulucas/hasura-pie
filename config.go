package pie

import (
	"github.com/joho/godotenv"
	"os"
	"strings"
)

func init() {
	logger := NewLogger()

	file := ".env"

	if IsProduction() {
		file += ".prod"
	} else {
		logger.WithField("core", "env").Warnf("Running at dev mode")
	}
	if err := godotenv.Load(file); err == nil {
		logger.WithField("core", "env").Infof("Load env file: %s", file)
		return
	}
}

func IsProduction() bool {
	return strings.ToLower(os.Getenv("APP_PRODUCTION")) == "true"
}
