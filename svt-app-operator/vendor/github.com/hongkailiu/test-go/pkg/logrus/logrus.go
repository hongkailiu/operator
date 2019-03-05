package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	//https://github.com/sirupsen/logrus/blob/master/example_basic_test.go
	var log = logrus.New()
	log.Formatter = new(logrus.TextFormatter)
	log.Out = os.Stdout
	log.Formatter.(*logrus.TextFormatter).DisableColors = true // remove colors
	//log.Formatter.(*logrus.TextFormatter).DisableTimestamp = true // remove timestamp from test output

	log.WithFields(logrus.Fields{
		"animal": "walrus",
	}).Info("A walrus appears")

}
