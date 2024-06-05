package log

import (
	"time"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
)

func Init() {
	logger = logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.Hooks.Add(lfshook.NewHook(
		getPathMap(
			"./log/wallet-system-error.log",
			"./log/wallet-system-info.log",
		),
		&logrus.TextFormatter{
			TimestampFormat: time.RFC3339,
			FullTimestamp:   true,
		},
	))
}

func getPathMap(errorFilePath string, infoFilePath string) lfshook.PathMap {
	return lfshook.PathMap{
		logrus.PanicLevel: errorFilePath,
		logrus.ErrorLevel: errorFilePath,
		logrus.FatalLevel: errorFilePath,
		logrus.WarnLevel:  errorFilePath,
		logrus.InfoLevel:  infoFilePath,
		logrus.DebugLevel: infoFilePath,
	}
}
