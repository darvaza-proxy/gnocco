package gnocco

import (
	"github.com/darvaza-proxy/gnocco/shared/cblog"
	"github.com/darvaza-proxy/gnocco/shared/log"
)

func (cf *Gnocco) Logger() log.Logger {
	return cf.logger
}

func newLogger(cf *Gnocco) log.Logger {
	logger := cblog.New()

	if cf.Log.Stdout {
		logger.SetLogger("console", nil)
	}

	if cf.Log.File != "" {
		cfg := map[string]interface{}{"file": cf.Log.File}
		logger.SetLogger("file", cfg)
		logger.Info("Logger started")
	}

	return logger
}
