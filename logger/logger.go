package logger

import (
	"os"

	"github.com/Sirupsen/logrus"
)

var standardLogger = logrus.StandardLogger()

func StandardLogger() *logrus.Logger {
	return standardLogger
}

func Enable() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)
}

func Debug(format string, args ...interface{}) {
	StandardLogger().Debugf(format, args...)
}

func Info(format string, args ...interface{}) {
	StandardLogger().Infof(format, args...)
}

func Warning(format string, args ...interface{}) {
	StandardLogger().Warningf(format, args...)
}

func Error(format string, args ...interface{}) {
	StandardLogger().Errorf(format, args...)
}
