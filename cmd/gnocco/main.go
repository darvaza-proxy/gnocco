// Gnocco is a little cache of goodness
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/darvaza-proxy/gnocco/server/gnocco"
	"github.com/darvaza-proxy/gnocco/shared/version"
)

func main() {
	var confFile string
	var vrs bool
	flag.StringVar(&confFile, "f", "", "specify the config file, if empty will try gnocco.conf and /etc/gnocco/gnocco.conf.")
	flag.BoolVar(&vrs, "v", false, "program version")
	flag.Parse()

	if vrs {
		fmt.Fprintf(os.Stdout, "Gnocco version %s, build date %s\n", version.Version, version.BuildDate)
		os.Exit(0)
	}

	cf, err := gnocco.NewFromFilename(confFile)
	if err != nil {
		panic(err)
	}

	logger := cf.Logger()

	aserver := &gnocco.GnoccoServer{
		Host:       cf.Listen.Host,
		Port:       cf.Listen.Port,
		MaxJobs:    cf.MaxJobs,
		MaxQueries: cf.MaxQueries,
	}

	aserver.Run()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan)

	for {
		sign := <-signalChan
		switch sign {
		case syscall.SIGTERM:
			aserver.ShutDown()
			logger.Fatal().Print("Got SIGTERM, stoping as requested")
		case syscall.SIGINT:
			aserver.ShutDown()
			logger.Fatal().Print("Got SIGINT, stoping as requested")
		case syscall.SIGUSR2:
			logger.Info().Print("Got SIGUSR2, dumping cache")
			aserver.DumpCache()
		case syscall.SIGURG:
		default:
			logger.Warn().Printf("I received %v signal", sign)
		}
	}
}
