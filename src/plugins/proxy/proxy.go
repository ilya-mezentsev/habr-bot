package proxy

import (
	"fmt"
	"net/http"
	"net/url"
	"sync/atomic"
)

var (
	checkedIpCount   int64
	availableIpFound bool
)

func GetAvailableIp(ips []string) (*url.URL, error) {
	ipCount := int64(len(ips))
	availableIps := make(chan string)

	for _, ip := range ips {
		go checkIp(ip, availableIps)
	}

	for {
		select {
		case ip := <-availableIps:
			availableIpFound = true
			return url.Parse(fmt.Sprintf("http://%s", ip))
		default:
			if checkedIpCount >= ipCount {
				return nil, noAvailableIps
			}
		}
	}
}

func increaseCheckedIpCount() {
	atomic.AddInt64(&checkedIpCount, 1)
}

func checkIp(ip string, availableIps chan<- string) {
	_, err := http.Get(fmt.Sprintf("http://%s", ip))
	if err == nil && !availableIpFound {
		availableIps <- ip
	}

	increaseCheckedIpCount()
}
