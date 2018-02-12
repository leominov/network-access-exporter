package main

import "net"

type Item struct {
	Resource string
	Host     string
	Port     int
}

func (i *Item) Lookup() ([]net.IP, error) {
	result := []net.IP{}
	ipAddresses, err := net.LookupIP(i.Host)
	if err != nil {
		return result, err
	}
	return ipAddresses, nil
}
