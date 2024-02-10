package agent

import (
	"fmt"
	"net"
)

func getLocalIp() (string, error) {
	var localIp string
	ifaces, err := net.Interfaces()
	if err != nil {
		return localIp, fmt.Errorf("error on getting local ip: %w", err)
	}

	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			if len(addr.String()) != 0 {
				return addr.String(), nil
			}
		}
	}

	return localIp, nil
}
