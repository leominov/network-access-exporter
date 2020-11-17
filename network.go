package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	minTCPPort = 0
	maxTCPPort = 65535
)

func IsTCPPortAvailable(ipAddress net.IP, port int, timeout time.Duration) bool {
	if port < minTCPPort || port > maxTCPPort {
		return false
	}
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ipAddress, port), timeout)
	if err != nil {
		return false
	}
	if err := conn.Close(); err != nil {
		return false
	}
	return true
}

func IsTCPPortAvailableIface(iface string, ipAddress net.IP, port int, timeout time.Duration) bool {
	ief, err := net.InterfaceByName(iface)
	if err != nil {
		return false
	}
	addrs, err := ief.Addrs()
	if err != nil {
		return false
	}
	ipAddr := addrs[0].(*net.IPNet).IP
	if IsIPv6(ipAddr.String()) {
		ipAddr = addrs[1].(*net.IPNet).IP
	}
	tcpAddr := &net.TCPAddr{
		IP: ipAddr,
	}
	d := net.Dialer{LocalAddr: tcpAddr, Timeout: timeout}
	_, err = d.Dial("tcp", fmt.Sprintf("%s:%d", ipAddress, port))
	if err != nil {
		return false
	}
	return true
}

func IsIPv6(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil && strings.Contains(str, ":")
}
