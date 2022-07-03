package algo

import (
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	p = uint64(2654435761)
)

func randInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}

func addr2Int(addr net.Addr) uint64 {
	ipAddrAndPort := strings.Split(addr.String(), ":")
	//x.x.x.x
	ipAddr := strings.Split(ipAddrAndPort[0], ".")
	first, _ := strconv.Atoi(ipAddr[0])
	second, _ := strconv.Atoi(ipAddr[1])
	third, _ := strconv.Atoi(ipAddr[2])
	fourth, _ := strconv.Atoi(ipAddr[3])

	return uint64(first)*16777216 + uint64(second)*65536 + uint64(third)*256 + uint64(fourth)
}

func IpHash(healthList []int, remoteAddr net.Addr, n int) int {
	ip := addr2Int(remoteAddr)
	k := int(ip * p % uint64(n))
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

func Rand(healthList []int) int {
	return healthList[randInt(0, len(healthList))]
}

func RoundRobin(healthList []int, round int) int {
	tmp := round + 1
	max := healthList[len(healthList)-1]
	if tmp > max {
		tmp = 0
	}

	for i := 0; i < len(healthList); i++ {
		if tmp == healthList[i] {
			return tmp
		}
		if tmp > healthList[i] {
			continue
		}
		if tmp < healthList[i] {
			return healthList[i]
		}
	}
	return -1
}
