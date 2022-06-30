package proxy

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/adrian-lin-1-0-0/gobalancer/pkg/logger"
)

type listeners []listener

type listener struct {
	Port              int
	SSL               bool
	SSLCertificate    string
	SSLCertificateKey string
	Instances         Instances
	Upstream          int
	Listener          net.Listener
	HealthCheck       int
	Round             int
	RWMux             sync.RWMutex
	Algo              string
}

type Listeners interface {
	Run() error
	Close() error
}

func (ls *listeners) Run() error {
	for i := 0; i < len(*ls); i++ {
		if err := (*ls)[i].init(); err != nil {
			return err
		}
	}
	return nil
}

func (ls *listeners) Close() error {
	for i := 0; i < len(*ls); i++ {
		if err := (*ls)[i].close(); err != nil {
			return err
		}
	}
	return nil
}

func (l *listener) close() error {
	return (*l).Listener.Close()
}

func (l *listener) init() error {
	address := fmt.Sprintf(":%d", l.Port)
	var nl net.Listener
	var err error
	if l.SSL {
		cert, err := tls.LoadX509KeyPair(l.SSLCertificate, l.SSLCertificateKey)
		if err != nil {
			return err
		}
		nl, err = tls.Listen("tcp", address, &tls.Config{Certificates: []tls.Certificate{cert}})
		if err != nil {
			return err
		}
	} else {
		nl, err = net.Listen("tcp", address)
		if err != nil {
			return err
		}
	}
	l.Listener = nl
	go l.run(&nl)

	return nil
}

func (l *listener) healthList() []int {
	list := []int{}
	for i := 0; i < len((*l).Instances); i++ {
		if (*l).Instances[i].isAlive() {
			list = append(list, i)
		}
	}
	return list
}

func (l *listener) run(nl *net.Listener) {
	defer (*nl).Close()
	go func() {
		for {
			l.Instances.healthcheck()
			time.Sleep(time.Duration(l.HealthCheck) * time.Second)
		}
	}()

	logger.Log.Info(fmt.Sprintf(":%d -> :%d %s", l.Port, l.Upstream, l.Algo))

	for {
		conn, err := (*nl).Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				logger.Log.Error(err.Error())
				continue
			}
			return
		}
		go l.handleConnection(conn)
	}
}

func (l *listener) handleConnection(conn net.Conn) {
	defer conn.Close()

	if idx := l.pickInscance(l.Algo, conn.RemoteAddr()); idx > -1 {
		logger.Log.Info(l.Port, "->", l.Instances[idx].Addr.String())
		upstream, err := net.Dial("tcp", l.Instances[idx].Addr.String())
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}
		defer upstream.Close()
		go func() {
			defer upstream.Close()
			defer conn.Close()
			io.Copy(upstream, conn)
		}()

		io.Copy(conn, upstream)
		return
	}
	logger.Log.Warn("No instance is alive")
}
