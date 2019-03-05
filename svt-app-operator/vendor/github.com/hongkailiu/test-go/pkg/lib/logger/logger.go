package logger

/*import (
	"os"

	"github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
	"log/syslog"
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
}*/

import (
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
)

// NewLogger returns a logrus logger which writes both on standard error and a file
// https://github.com/Sirupsen/logrus/issues/230
// if we do not need the color on terminal, the following will work as described here
// https://github.com/Sirupsen/logrus/issues/230#issuecomment-381639138
// mw := io.MultiWriter(os.Stdout, logFile)
// logrus.SetOutput(mw)
func NewLogger(logFilePath string) *logrus.Logger {
	if logger != nil {
		return logger
	}

	pathMap := lfshook.PathMap{
		logrus.InfoLevel:  logFilePath,
		logrus.ErrorLevel: logFilePath,
	}

	logger = logrus.New()

	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	logger.Hooks.Add(lfshook.NewHook(
		pathMap,
		//&logrus.JSONFormatter{},
		&logrus.TextFormatter{},
	))
	return logger
}
