package pie

import (
	"github.com/sirupsen/logrus"
)

type Logger interface {
	logrus.FieldLogger
}

func NewLogger() Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	return logger
}
