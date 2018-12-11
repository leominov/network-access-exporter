package main

import "net"

type Item struct {
	Resource string `yaml:"resource"`
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

func (i *Item) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var stringType string
	if err := unmarshal(&stringType); err == nil {
		*i = Item{
			Resource: stringType,
			Group:    "",
		}
		return nil
	}
	var objType Item
	if err := unmarshal(&objType); err != nil {
		return err
	}
	*i = objType
	return nil
}
