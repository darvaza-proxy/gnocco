// Package main is the main executable
package main

import (
	"flag"

	"darvaza.org/gnocco"
	"darvaza.org/gnocco/zerolog"
)

func main() {
	var confFile string

	flag.StringVar(&confFile, "f", "", "specify the config file, "+
		"if empty will try gnocco.conf and /etc/gnocco/gnocco.conf.")
	flag.Parse()

	cf := new(gnocco.Config)
	if confFile == "" {
		confFile = gnocco.DefaultConfigFile
	}
	err := cf.ReadInFile(confFile)
	if err != nil {
		panic(err)
	}
	logger := zerolog.NewLogger(gnocco.DefaultLogLevel)
	if err := gnocco.Run(cf, logger); err != nil {
		logger.Fatal().Print(err)
		return
	}
}
