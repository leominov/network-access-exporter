package main

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

type Item struct {
	Resource string `yaml:"addr"`
	Group    string `yaml:"group,omitempty"`
	Network  string `yaml:"network,omitempty"`
	Host     string `yaml:"-"`
	Port     int    `yaml:"-"`
	Iface    string `yaml:"iface"`
}

func ParseResource(resource string, iface_optional ...string) (*Item, error) {
	iface := ""
	if len(iface_optional) > 0 {
		iface = iface_optional[0]
	}
	network := getResourceNetwork(resource)
	resource = strings.TrimPrefix(resource, fmt.Sprintf("%s://", network))
	host, port, err := net.SplitHostPort(resource)
	if err != nil {
		return nil, err
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("incorrent port in item: %+v", resource)
	}
	return &Item{
		Network:  network,
		Host:     host,
		Port:     portInt,
		Iface:    iface,
		Resource: resource,
	}, nil
}

func (i *Item) Lookup() ([]net.IP, error) {
	result := []net.IP{}
	ipAddresses, err := net.LookupIP(i.Host)
	if err != nil {
		return result, err
	}
	return ipAddresses, nil
}

func getResourceNetwork(resource string) string {
	if !strings.Contains(resource, "://") {
		return "tcp"
	}
	u, err := url.Parse(resource)
	if err != nil {
		return "tcp"
	}
	return u.Scheme
}
