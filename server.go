package gnocco

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"darvaza.org/slog"
	"github.com/miekg/dns"
)

const (
	// DefaultConfigFile specifies the default filename of the config file
	DefaultConfigFile = "gnocco.conf"
	// DefaultLogLevel specifies the log level we handle by default
	DefaultLogLevel = slog.Debug
)

var (
	cfg     Config
	cfgFile string
	log     slog.Logger
)

// Run creates and runs a gnocco server using the given config
func Run(cf *Config, logger slog.Logger) error {
	log = logger
	addr := net.JoinHostPort(cf.Listen.Host, cf.Listen.Port)

	h, err := NewHandler("")
	if err != nil {
		return err
	}
	dns.Handle(".", h)

	udpServer := &dns.Server{Addr: addr, Net: "udp"}
	tcpServer := &dns.Server{Addr: addr, Net: "tcp"}

	errChan := make(chan error)

	go runServer(udpServer, errChan)
	go runServer(tcpServer, errChan)

	sigHandlers()

	err = shutdown(udpServer)
	if err != nil {
		errChan <- err
	}
	err = shutdown(tcpServer)
	if err != nil {
		errChan <- err
	}

	if len(errChan) > 0 {
		return <-errChan
	}
	return nil
}

func runServer(srv *dns.Server, errChan chan error) {
	err := srv.ListenAndServe()
	if err != nil {
		errChan <- err
	}
}

func shutdown(srv *dns.Server) error {
	return srv.Shutdown()
}

func sigHandlers() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals)
	for {
		sig := <-signals
		switch sig {
		case syscall.SIGTERM, syscall.SIGINT:
			log.Info().Printf("Got %s, doing as requested\n", sig.String())
			return
		case syscall.SIGUSR2:
			log.Info().Printf("Got SIGUSR2, dumping cache")
		case syscall.SIGURG:
		default:
			log.Info().Printf("Received \"%v\" which is not registered\n", sig)
		}
	}
}
