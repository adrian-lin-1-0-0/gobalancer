package proxy

import (
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

func randInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}

func roundRobin(l *listener, healthList []int) int {
	l.RWMux.Lock()
	defer l.RWMux.Unlock()
	tmp := l.Round + 1
	max := healthList[len(healthList)-1]
	if tmp > max {
		tmp = 0
	}

	for i := 0; i < len(healthList); i++ {
		if tmp == healthList[i] {
			l.Round = tmp
			return tmp
		}
		if tmp > healthList[i] {
			continue
		}
		if tmp < healthList[i] {
			l.Round = healthList[i]
			return healthList[i]
		}
	}
	return -1
}

func addr2Int(addr net.Addr) int {
	ipAddrAndPort := strings.Split(addr.String(), ":")
	//x.x.x.x
	ipAddr := strings.Split(ipAddrAndPort[0], ".")
	first, _ := strconv.Atoi(ipAddr[0])
	second, _ := strconv.Atoi(ipAddr[1])
	third, _ := strconv.Atoi(ipAddr[2])
	fourth, _ := strconv.Atoi(ipAddr[3])

	return first*16777216 + second*65536 + third*256 + fourth
}

func ipHash(remoteAddr net.Addr, healthList []int, n int) int {
	ip := addr2Int(remoteAddr)
	k := ip * 2654435761 % n
	if k > healthList[len(healthList)-1] {
		return healthList[0]
	}

	for i := 0; i < len(healthList); i++ {
		if k == healthList[i] {
			return k
		}
		if k < healthList[i] {
			return healthList[i]
		}
	}
	return k
}

func (l *listener) pickInscance(al string, remoteAddr net.Addr) int {
	healthList := l.healthList()
	if len(healthList) == 0 {
		return -1
	}

	switch al {
	case "rand":
		return healthList[randInt(0, len(healthList))]
	case "round-robin":
		return roundRobin(l, healthList)
	case "ip-hash":
		return ipHash(remoteAddr, healthList, len(l.Instances))
	default:
		return healthList[randInt(0, len(healthList))]
	}
}
