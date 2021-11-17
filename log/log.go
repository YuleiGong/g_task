package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

var glog *logrus.Logger

func init() {
	glog = logrus.New()
	glog.SetFormatter(&logrus.JSONFormatter{})
	glog.SetOutput(os.Stdout)
	glog.SetLevel(logrus.InfoLevel)
	glog.Out = os.Stdout
}

func Info(format string, args ...interface{}) {
	glog.Infof(format, args...)
}

func Error(format string, args ...interface{}) {
	glog.Errorf(format, args...)
}
