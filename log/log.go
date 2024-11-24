package logs

import (
	"github.com/sirupsen/logrus"
)

func NewLogDebug() logrus.FieldLogger {
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	l.SetFormatter(&logrus.TextFormatter{})
	l.SetReportCaller(false)
	l.Trace("Debug log")
	return l
}
func NewLogInfo() logrus.FieldLogger {
	l := logrus.New()
	l.SetLevel(logrus.InfoLevel)
	l.SetFormatter(&logrus.TextFormatter{})
	l.SetReportCaller(true)
	l.Trace("Info log")
	return l
}
