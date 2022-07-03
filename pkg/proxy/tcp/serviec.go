package tcp

import (
	"context"
	"fmt"
	"net"

	"github.com/adrian-lin-1-0-0/gobalancer/pkg/config"
)

func Init(config_ config.Config, ctx context.Context, logger Logger) (Listeners, error) {
	serviceConfig := config_.Get()
	ls := NewListeners()
	for i := 0; i < len(serviceConfig.Listener); i++ {
		l, err := initListener(&serviceConfig.Listener[i], &serviceConfig.Instance, ctx, logger)
		if err != nil {
			return nil, err
		}
		if l == nil {
			continue
		}
		ls.Add((l).Get())
	}

	if ls.Len() == 0 {
		return nil, nil
	}

	return ls, nil
}

func initListener(listenerConfig *config.Listener, instancesConfig *[]config.Instance, ctx context.Context, logger Logger) (Listener, error) {
	if listenerConfig.Protocol == "udp" {
		return nil, nil
	}
	l := &listener{
		ctx:                 ctx,
		Logger:              logger,
		Port:                listenerConfig.Port,
		SSL:                 listenerConfig.SSL,
		SSLCertificate:      listenerConfig.SSLCertificate,
		SSLCertificateKey:   listenerConfig.SSLCertificateKey,
		HealthCheckInterval: listenerConfig.HealthCheckInterval,
		Nagle:               listenerConfig.Nagle,
	}

	if listenerConfig.HealthCheckInterval == 0 {
		l.HealthCheckInterval = 30
	} else {
		l.HealthCheckInterval = listenerConfig.HealthCheckInterval
	}

	if listenerConfig.Algo == "" {
		l.Algo = "rand"
	} else {
		l.Algo = listenerConfig.Algo
	}
	logger.Info(fmt.Sprintf(":%d -> :%d (%s)", listenerConfig.Port, listenerConfig.Upstream, l.Algo))
	for i := 0; i < len(*instancesConfig); i++ {
		tcpAddr, err := net.ResolveTCPAddr(
			"tcp",
			fmt.Sprintf("%s:%d", (*instancesConfig)[i].Addr, listenerConfig.Upstream),
		)
		if err != nil {
			return nil, err
		}
		l.Instances.Add(NewInstance(*tcpAddr))
	}

	return l, nil
}
