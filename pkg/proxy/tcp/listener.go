package tcp

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/adrian-lin-1-0-0/gobalancer/pkg/proxy/algo"
)

type listener struct {
	Round       int
	RWMux       sync.RWMutex
	NetListener net.Listener
	ctx         context.Context

	//setting
	Logger              Logger
	Port                int
	SSL                 bool
	SSLCertificate      string
	SSLCertificateKey   string
	Instances           Instances
	HealthCheckInterval int
	Algo                string
	Nagle               bool
}

type Listener interface {
	Get() *listener
}

type listeners []*listener

type Listeners interface {
	Len() int
	Run() error
	Add(listener *listener)
}

func (listeners_ *listeners) Len() int {
	return len(*listeners_)
}

func (listeners_ *listeners) Run() error {
	for i := 0; i < len(*listeners_); i++ {
		if err := (*listeners_)[i].init(); err != nil {
			return err
		}
	}
	return nil
}

func (listeners_ *listeners) Add(listener *listener) {
	*listeners_ = append(*listeners_, listener)
}

func NewListeners() Listeners {
	return &listeners{}
}

func NewListener() *listener {
	return &listener{}
}

func (listener_ *listener) Get() *listener {
	return listener_
}

func (listener_ *listener) init() error {
	address := fmt.Sprintf(":%d", listener_.Port)
	var err error
	if listener_.SSL {
		cert, err := tls.LoadX509KeyPair(listener_.SSLCertificate, listener_.SSLCertificateKey)
		if err != nil {
			return err
		}
		listener_.NetListener, err = tls.Listen("tcp", address, &tls.Config{Certificates: []tls.Certificate{cert}})
		if err != nil {
			return err
		}
	} else {
		listener_.NetListener, err = net.Listen("tcp", address)
		if err != nil {
			return err
		}
	}
	go listener_.run()
	return nil
}

func (listener_ *listener) run() {
	defer listener_.NetListener.Close()
	go listener_.HealthCheck()

	for {
		select {
		case <-listener_.ctx.Done():
			return
		default:
			break
		}
		conn, err := listener_.NetListener.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				listener_.Logger.Error(err.Error())
				continue
			}
			return
		}
		go listener_.handleConnection(conn)
	}

}

func (listener_ *listener) handleConnection(conn net.Conn) {
	defer conn.Close()
	if idx := listener_.pickInstance(conn.RemoteAddr()); idx > -1 {
		listener_.Logger.Info(conn.RemoteAddr(), "->", listener_.Instances[idx].Addr.String())
		upstream, err := net.DialTCP("tcp", nil, &listener_.Instances[idx].Addr)
		if err != nil {
			listener_.Logger.Error(err.Error())
			return
		}
		defer upstream.Close()
		if !listener_.Nagle {
			upstream.SetNoDelay(true)
		}
		go func() {
			defer conn.Close()
			io.Copy(upstream, conn)
		}()
		io.Copy(conn, upstream)
		return
	}
	listener_.Logger.Warn("No instance is alive")
}

func (listener_ *listener) healthList() []int {
	list := []int{}
	for i := 0; i < len((*listener_).Instances); i++ {
		if (*listener_).Instances[i].isAlive() {
			list = append(list, i)
		}
	}
	return list
}

func (listener_ *listener) getRound() int {
	listener_.RWMux.RLock()
	defer listener_.RWMux.RUnlock()
	return listener_.Round
}

func (listener_ *listener) setRound(round int) {
	listener_.RWMux.Lock()
	defer listener_.RWMux.Unlock()
	listener_.Round = round
}

func (listener_ *listener) pickInstance(remoteAddr net.Addr) int {
	healthList := listener_.healthList()
	if len(healthList) == 0 {
		return -1
	}
	switch listener_.Algo {
	case "rand":
		return algo.Rand(healthList)
	case "round-robin":
		round := listener_.getRound()
		if newRound := algo.RoundRobin(healthList, round); newRound != -1 {
			listener_.setRound(newRound)
			return newRound
		}
		return -1
	case "ip-hadh":
		return algo.IpHash(healthList, remoteAddr, len(listener_.Instances))
	default:
		return algo.Rand(healthList)
	}
}

func (listener_ *listener) close() error {
	return listener_.NetListener.Close()
}

func (listener_ *listener) HealthCheck() {
	select {
	case <-listener_.ctx.Done():
		return
	default:
		listener_.Instances.healthcheck()
		time.Sleep(time.Duration(listener_.HealthCheckInterval) * time.Second)
	}
}
