package logger

import (
	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/sirupsen/logrus"
	"os"
)

// SetupLogger is function to setup logger init
func SetupLogger() {
	childFormatter := new(logrus.TextFormatter)
	childFormatter.TimestampFormat = "2006-01-02 15:04:05"
	childFormatter.FullTimestamp = true

	formatter := runtime.Formatter{
		ChildFormatter: childFormatter,
		Line:           true,
		File:           true,
	}

	logrus.SetFormatter(&formatter)
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}
