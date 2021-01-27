package config

import (
	"errors"
	"fmt"
	"net"
)

// ReadSystemIP will return the ip address
func ReadSystemIP(netint string) (string, error) {
	iface, err := net.InterfaceByName(netint)
	if err != nil {
		return "", err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		returnIP := fmt.Sprintf("%+v", ip)
		return returnIP, nil
	}

	return "", errors.New("Interface was not found")
}
