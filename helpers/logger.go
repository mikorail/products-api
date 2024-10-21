package helpers

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func InitLogger() {
	Logger.Out = os.Stdout
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	Logger.SetLevel(logrus.DebugLevel)
}

func Info(message string, fields logrus.Fields) {
	Logger.WithFields(fields).Info(message)
}

func Warn(message string, fields logrus.Fields) {
	Logger.WithFields(fields).Warn(message)
}

func Error(message string, fields logrus.Fields) {
	Logger.WithFields(fields).Error(message)
}

func Debug(message string, fields logrus.Fields) {
	Logger.WithFields(fields).Debug(message)
}
