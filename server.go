package gnocco

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/miekg/dns"
	"github.com/spf13/viper"
)

// ListenAndServe creates a new Gnocco server and starts listening
func ListenAndServe() error {
	NewConfig()

	addr := net.JoinHostPort(viper.GetString("Listen.Host"), viper.GetString("Listen.Port"))

	dns.HandleFunc(".", HandleRequest)

	udpServer := &dns.Server{Addr: addr, Net: "udp"}
	tcpServer := &dns.Server{Addr: addr, Net: "tcp"}

	errChan := make(chan error)

	go func() {
		err := udpServer.ListenAndServe()
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		err := tcpServer.ListenAndServe()
		if err != nil {
			errChan <- err
		}
	}()

	sigHandlers()

	err := shutdown(udpServer)
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
			fmt.Printf("Got %s, doing as requested\n", sig.String())
			return
		case syscall.SIGUSR2:
			fmt.Println("Got SIGUSR2, dumping cache")
		case syscall.SIGURG:
		default:
			fmt.Printf("Received signal %v\n", sig)
		}
	}
}
