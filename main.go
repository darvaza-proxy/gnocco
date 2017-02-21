package main

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	logger *GnoccoLogger
)

func main() {
	cf := loadConfig()
	initLogger()

	server := &Server{
		host:  cf.Listen.Host,
		port:  cf.Listen.Port,
		dotcp: cf.DoTCP,
		user:  cf.User,
		group: cf.Group,
	}

	server.Run()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan)

	for {
		sign := <-signalChan
		switch sign {
		case syscall.SIGTERM:
			logger.Fatal("Got SIGTERM, stoping as requested")
		case syscall.SIGINT:
			logger.Fatal("Got SIGINT, stoping as requested")
		case syscall.SIGUSR2:
			logger.Info("Got SIGUSR2, dumping cache to stdout")
			server.DumpCache()
		default:
			logger.Warn("I received %s signal", sign)
		}
	}
}

func initLogger() {
	logger = NewLogger()

	if Config.Log.Stdout {
		logger.SetLogger("console", nil)
	}

	if Config.Log.File != "" {
		cfg := map[string]interface{}{"file": Config.Log.File}
		logger.SetLogger("file", cfg)
		logger.Info("Logger started")
	}

}
