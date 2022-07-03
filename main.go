package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/adrian-lin-1-0-0/gobalancer/pkg/config"
	"github.com/adrian-lin-1-0-0/gobalancer/pkg/logger"
	"github.com/adrian-lin-1-0-0/gobalancer/pkg/proxy/tcp"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	serviceConfig, err := config.Load()
	if err != nil {
		log.Fatalf("err : %s", err)
	}
	tcpService, err := tcp.Init(serviceConfig, ctx, logger.Log)
	if err != nil {
		log.Fatalf("err : %s", err)
	}

	logger.Log.Info("service startting...")

	if tcpService != nil {
		tcpService.Run()
	}
	clean(cancel)
}

func clean(cancel context.CancelFunc) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	s := <-signalChan
	cancel()
	logger.Log.Info(fmt.Sprintf("Got signal : %s ,Stoppping services ...", s.String()))

}
