package main

import (
	"fmt"
	"os"

	"github.com/hongkailiu/test-go/pkg/http"
	"github.com/hongkailiu/test-go/pkg/http/info"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("test-go", "A web application for test.").Version(info.VERSION)

	debug = app.Flag("debug", "Enable debug mode.").Short('v').Default("false").Bool()

	get          = app.Command("get", "get resources.")
	getResource  = get.Arg("resource", "targeting resource.").Required().String()
	secretLength = get.Flag("secretLength", "the length of the secret.").Short('l').Default("32").Int()

	start = app.Command("start", "Start server")
)

func init() {

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)
}

func main() {
	command := kingpin.MustParse(app.Parse(os.Args[1:]))
	if debug != nil && *debug {
		log.SetLevel(log.DebugLevel)
	}
	switch command {
	case get.FullCommand():
		log.Info(get.FullCommand())
		switch *getResource {
		case "secret":
			fmt.Println(http.GetSecret(*secretLength))

		}
	case start.FullCommand():
		log.Info(start.FullCommand())
		http.Run()
	}

}
