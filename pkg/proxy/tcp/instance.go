package tcp

import (
	"net"
	"sync"
	"time"
)

const (
	//400ms, 0.4 seconds
	pingTimeOut = time.Millisecond * 400
)

type Instance struct {
	RWMux sync.RWMutex
	Alive bool
	Addr  net.TCPAddr
}

type Instances []*Instance

func (instances *Instances) Add(instance *Instance) {
	*instances = append(*instances, instance)
}

func NewInstance(addr net.TCPAddr) *Instance {
	return &Instance{
		Addr: addr,
	}
}

func (instances *Instances) healthcheck() {
	for i := 0; i < len(*instances); i++ {
		(*instances)[i].heathcheck()
	}
}

func (instance *Instance) isAlive() bool {
	instance.RWMux.RLock()
	defer instance.RWMux.RUnlock()
	alive := instance.Alive
	return alive
}

func (instance *Instance) heathcheck() {
	isAlive := healthcheck(instance.Addr)
	if isAlive == instance.isAlive() {
		return
	}
	instance.RWMux.Lock()
	defer instance.RWMux.Unlock()
	instance.Alive = isAlive
	return
}

func healthcheck(addr net.TCPAddr) bool {
	conn, err := net.DialTimeout("tcp", addr.String(), pingTimeOut)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
