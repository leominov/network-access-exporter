package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type Exporter struct {
	config *Config
}

func NewExporter(config *Config) *Exporter {
	return &Exporter{
		config: config,
	}
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	var avFloat float64
	var isIPV6Address string
	for _, item := range e.config.Items {
		ipAddresses, err := item.Lookup()
		if err != nil {
			logrus.Errorf("Cant get IP address: %v", err)
			ch <- prometheus.MustNewConstMetric(loopkupError, prometheus.GaugeValue, 1.0, item.Resource)
			continue
		}
		for _, ipAddress := range ipAddresses {
			logrus.Debugf("Checking %s with port %d on %s...", item.Host, item.Port, ipAddress.String())
			if IsTCPPortAvailable(ipAddress, item.Port, e.config.ConnectionTimeout) {
				avFloat = 1.0
			} else {
				logrus.Warnf("TCP port not available: %s on %s", item.Resource, ipAddress.String())
				avFloat = 0.0
			}
			if ok := IsIPv6(ipAddress.String()); ok {
				isIPV6Address = "1"
			} else {
				isIPV6Address = "0"
			}
			ch <- prometheus.MustNewConstMetric(allowedResource, prometheus.GaugeValue, avFloat, item.Resource, ipAddress.String(), isIPV6Address)
		}
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- allowedResource
	ch <- loopkupError
}
