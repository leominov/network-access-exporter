package main

import (
	"net"
)

type Item struct {
	Resource string `yaml:"addr"`
	Group    string `yaml:"group,omitempty"`
	Host     string `yaml:"-"`
	Port     int    `yaml:"-"`
}

func (i *Item) Lookup() ([]net.IP, error) {
	result := []net.IP{}
	ipAddresses, err := net.LookupIP(i.Host)
	if err != nil {
		return result, err
	}
	return ipAddresses, nil
}
