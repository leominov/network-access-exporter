package main

import (
	"time"

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
	var (
		avFloat           float64
		isIPV6Address     string
		startTime         time.Time
		durationInSeconds float64
	)
	for _, item := range e.config.Items {
		startTime = time.Now()
		ipAddresses, err := item.Lookup()
		if err != nil {
			logrus.Errorf("Cant get IP address: %v", err)
			ch <- prometheus.MustNewConstMetric(loopkupError, prometheus.GaugeValue, 1.0, item.Resource)
			continue
		}
		durationInSeconds = time.Now().Sub(startTime).Seconds()
		ch <- prometheus.MustNewConstMetric(loopkupDuration, prometheus.GaugeValue, durationInSeconds, item.Resource)
		for _, ipAddress := range ipAddresses {
			logrus.Debugf("Checking %s with port %d on %s...", item.Host, item.Port, ipAddress.String())
			startTime = time.Now()
			if ok := IsIPv6(ipAddress.String()); ok {
				isIPV6Address = "1"
			} else {
				isIPV6Address = "0"
			}
			if IsTCPPortAvailable(ipAddress, item.Port, e.config.ConnectionTimeout) {
				avFloat = 1.0
				durationInSeconds = time.Now().Sub(startTime).Seconds()
				ch <- prometheus.MustNewConstMetric(dialDuration, prometheus.GaugeValue, durationInSeconds, item.Resource, ipAddress.String(), isIPV6Address)
			} else {
				logrus.Warnf("TCP port not available: %s on %s", item.Resource, ipAddress.String())
				avFloat = 0.0
			}
			ch <- prometheus.MustNewConstMetric(allowedResource, prometheus.GaugeValue, avFloat, item.Resource, ipAddress.String(), isIPV6Address)
		}
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- allowedResource
	ch <- loopkupError
	ch <- loopkupDuration
	ch <- dialDuration
}
