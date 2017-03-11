package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/naoina/toml"
)

type cfg struct {
	User           string
	Group          string
	RootsFile      string
	PermissionsDir string
	Daemon         bool
	IterateResolv  bool

	Listen struct {
		Host string
		Port int
	}

	MaxJobs    int
	MaxQueries int

	Log struct {
		Stdout bool
		File   string
	}

	Cache struct {
		DumpInterval int
		Expire       int
		MaxCount     int
		CachePath    string
	}

	Hosts HostsCfg
}

type HostsCfg struct {
	Enable           bool
	Hosts_File       string
	Refresh_Interval uint
}

var (
	Config   cfg
	confFile string
)

func loadConfig() cfg {
	flag.StringVar(&confFile, "f", "/etc/gnocco/gnocco.conf", "specify the config file, defaults to /etc/gnocco/gnocco.conf.")

	flag.Parse()

	file, err := os.Open(confFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		logger.Fatal("Error %s occured.", err)
	}

	if err := toml.Unmarshal(buf, &Config); err != nil {
		logger.Fatal("Error %s occured.", err)
	}
	return Config
}
