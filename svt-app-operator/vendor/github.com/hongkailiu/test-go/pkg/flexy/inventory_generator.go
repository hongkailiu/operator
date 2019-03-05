package flexy

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Generate the inventory files
func Generate(inputPath string, config Config, outputFolder string) error {
	bytes, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	fmt.Println("==========\n" + string(bytes))
	files, err := ioutil.ReadDir(inputPath)
	if err != nil {
		return err
	}

	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		os.MkdirAll(outputFolder, os.ModePerm)
	}

	for _, f := range files {
		log.WithFields(log.Fields{"f.Name()": f.Name()}).Debug("a file found")
		if err != nil {
			return err
		}
		if err = generate(f.Name(), filepath.Join(inputPath, f.Name()), config, outputFolder); err != nil {
			return err
		}

	}

	return nil
}

func generate(name string, file string, config Config, outputFolder string) error {

	t, err := template.New(name).ParseFiles(file)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(filepath.Join(outputFolder, name))
	if err != nil {
		return err
	}
	err = t.Execute(outputFile, config)
	if err != nil {
		return err
	}

	return nil
}
