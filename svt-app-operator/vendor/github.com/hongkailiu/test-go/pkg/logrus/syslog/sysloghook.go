package syslog

import (
	"log/syslog"
	"os"

	"github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
)

// Logger with syslog hook
var Logger = logrus.New()

func init() {
	Logger.Formatter = &logrus.TextFormatter{FullTimestamp: true}
	Logger.Out = os.Stdout
	Logger.Level = logrus.DebugLevel

	hook, err := lSyslog.NewSyslogHook("", "", syslog.LOG_INFO, "")

	if err != nil {
		panic(err.Error())
	}

	Logger.Hooks.Add(hook)
}
