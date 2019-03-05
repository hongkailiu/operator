package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hongkailiu/test-go/pkg/lib/logger"
	"github.com/hongkailiu/test-go/pkg/lib/util"
	"github.com/hongkailiu/test-go/pkg/ocpsanity"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("ocp-sanity", "A script to check sanity of OCP installation.")
	_   = app.Version(ocpsanity.VERSION).HelpFlag.Short('h')
	logFilePath = ocpsanity.LogFilePath
	log = logger.NewLogger(logFilePath)
)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))
	log.WithFields(logrus.Fields{"logFilePath": logFilePath}).Info("logging to file")
	log.WithFields(logrus.Fields{"startTime": time.Now().Format(time.RFC3339)}).Info("Starting OCP sanity check")
	err := sanityCheck()
	log.WithFields(logrus.Fields{"endTime": time.Now().Format(time.RFC3339)}).Info("Finishing OCP sanity check")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
func sanityCheck() error {
	return ocpsanity.StartSanityCheck(getConfigPath())
}

func getConfigPath() string {
	configPath := os.Getenv("KUBECONFIG")
	if configPath != "" {
		return configPath
	}

	if home := util.HomeDir(); home != "" {
		return filepath.Join(home, ".kube", "config")
	}

	return ""
}
