package main

import (
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func initConfig() Config {
	f, err := os.Open("config.yaml")
	if err != nil {
		logrus.Error(err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		logrus.Error(err)
	}

	return cfg
}

// ParseEventHandler returns the respective handler object specified in the config file.
func ParseEventHandler(conf Config) Handler {

	var eventHandler Handler
	switch {
	case len(conf.Handler.Name) > 0:
		eventHandler = new(Sample)
	default:
		eventHandler = new(Default)
	}
	if err := eventHandler.Init(conf); err != nil {
		log.Fatal(err)
	}
	return eventHandler
}
