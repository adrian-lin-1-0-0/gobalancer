package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func init() {

	levelMap := map[string]logrus.Level{
		"trace": logrus.TraceLevel,
		"debug": logrus.DebugLevel,
		"info":  logrus.InfoLevel,
		"warn":  logrus.WarnLevel,
		"fatal": logrus.FatalLevel,
		"panic": logrus.PanicLevel,
	}

	if os.Getenv("LOGGER_JSON_FORMATTER") == "true" {
		Log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	var level string
	if os.Getenv("LOGGER_LEVEL") == "" {
		level = "info"
	} else {
		level = os.Getenv("LOGGER_LEVEL")
	}

	Log.SetLevel(levelMap[level])
	Log.SetReportCaller(false)
}
