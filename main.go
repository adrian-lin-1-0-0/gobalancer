package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/adrian-lin-1-0-0/gobalancer/pkg/config"
	"github.com/adrian-lin-1-0-0/gobalancer/pkg/logger"
	"github.com/adrian-lin-1-0-0/gobalancer/pkg/proxy"
)

func main() {
	serviceConfig, err := config.Load()
	if err != nil {
		log.Fatalf("err : %s", err)
	}
	service, err := proxy.InitService(serviceConfig)
	if err != nil {
		log.Fatalf("err : %s", err)
	}
	logger.Log.Info("Startting services ...")
	err = service.Run()
	if err != nil {
		log.Fatalf("err : %s", err)
	}
	clean(&service)
}

func clean(service *proxy.Listeners) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	s := <-signalChan
	logger.Log.Info(fmt.Sprintf("Got signal : %s ,Stoppping services ...", s.String()))
	if err := (*service).Close(); err != nil {
		logger.Log.Error(err.Error())
	}
}
