package agent

import (
	"fmt"
	"net"
)

func getLocalIP() (string, error) {
	var localIP string
	ifaces, err := net.Interfaces()
	if err != nil {
		return localIP, fmt.Errorf("error on getting local ip: %w", err)
	}

	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			if len(addr.String()) != 0 {
				return addr.String(), nil
			}
		}
	}

	return localIP, nil
}
