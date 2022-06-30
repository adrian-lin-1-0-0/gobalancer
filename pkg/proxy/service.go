package proxy

import (
	"fmt"
	"net"

	"github.com/adrian-lin-1-0-0/gobalancer/pkg/config"
)

func InitService(c config.Config) (Listeners, error) {
	sc := c.Get()

	ls := &listeners{}
	for i := 0; i < len(sc.Listener); i++ {

		hc := sc.Listener[i].HealthCheck
		if hc == 0 {
			hc = 30
		}

		algo := sc.Listener[i].Algo
		if algo == "" {
			algo = "rand"
		}

		l := listener{
			Port:              sc.Listener[i].Port,
			SSL:               sc.Listener[i].SSL,
			SSLCertificate:    sc.Listener[i].SSLCertificate,
			SSLCertificateKey: sc.Listener[i].SSLCertificateKey,
			Upstream:          sc.Listener[i].Upstream,
			HealthCheck:       hc,
			Algo:              algo,
		}
		for j := 0; j < len(sc.Instance); j++ {
			tcpAddr, err := net.ResolveTCPAddr(
				"tcp",
				fmt.Sprintf("%s:%d", sc.Instance[j].Addr, sc.Listener[i].Upstream))
			if err != nil {
				return nil, err
			}
			l.Instances = append(l.Instances, Instance{
				Addr: *tcpAddr,
			})
		}
		*ls = append(*ls, l)
	}
	return ls, nil
}
