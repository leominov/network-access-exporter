package main

import (
	"net"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

const (
	minTCPPort = 0
	maxTCPPort = 65535
)

type Exporter struct {
	config *Config
}

func IsTCPPortAvailable(item *Item) bool {
	if item.Port < minTCPPort || item.Port > maxTCPPort {
		return false
	}
	conn, err := net.DialTimeout("tcp", item.Resource, 100*time.Millisecond)
	if err != nil {
		return false
	}
	if err := conn.Close(); err != nil {
		return false
	}
	return true
}

func NewExporter(config *Config) *Exporter {
	return &Exporter{
		config: config,
	}
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	var avFloat float64
	for _, item := range e.config.Items {
		if IsTCPPortAvailable(item) {
			avFloat = 1.0
		} else {
			logrus.Warnf("TCP port not available: %s", item.Resource)
			avFloat = 0.0
		}
		ch <- prometheus.MustNewConstMetric(
			allowedResource, prometheus.GaugeValue, avFloat, item.Resource, item.Alias,
		)
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- allowedResource
}
