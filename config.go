package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	defaultsListenAddr        = ":9407"
	defaultsMetricsPath       = "/metrics"
	defaultsConnectionTimeout = 500 * time.Millisecond
)

var (
	connectionTimeout   = flag.Duration("timeout", defaultsConnectionTimeout, "Connection timeout")
	logLevel            = flag.String("log-level", "info", "Logging level")
	listenAddress       = flag.String("web.listen-address", defaultsListenAddr, "Listen address")
	metricsPath         = flag.String("web.telemetry-path", defaultsMetricsPath, "Metrics path")
	resourcesCollection = flag.String("resources", "", "Resources list")
)

type Config struct {
	ConnectionTimeout time.Duration
	LogLevel          string
	ListenAddr        string
	MetricsPath       string
	Items             []Item
}

type Item struct {
	Alias    string
	Resource string
	Host     string
	Port     int
}

func LoadConfig() (*Config, error) {
	nc := &Config{}
	nc.SetEmptyToDefaults()
	if err := nc.LoadFromFlags(); err != nil {
		return nil, err
	}
	if err := parseConfig(nc); err != nil {
		return nil, err
	}
	return nc, nil
}

func (c *Config) SetEmptyToDefaults() {
	if len(c.ListenAddr) == 0 {
		c.ListenAddr = defaultsListenAddr
	}
	if len(c.MetricsPath) == 0 {
		c.MetricsPath = defaultsMetricsPath
	}
	if c.ConnectionTimeout.Nanoseconds() == 0 {
		c.ConnectionTimeout = defaultsConnectionTimeout
	}
}

func (c *Config) LoadFromFlags() error {
	if connectionTimeout.Nanoseconds() != 0 {
		c.ConnectionTimeout = *connectionTimeout
	}
	if len(*logLevel) != 0 {
		c.LogLevel = *logLevel
	}
	if len(*listenAddress) != 0 {
		c.ListenAddr = *listenAddress
	}
	if len(*metricsPath) != 0 {
		c.MetricsPath = *metricsPath
	}
	if len(*resourcesCollection) != 0 {
		resources := strings.Split(*resourcesCollection, ",")
		for _, resourceRaw := range resources {
			resourceRaw = strings.TrimSpace(resourceRaw)
			if len(resourceRaw) == 0 {
				continue
			}
			resource := strings.Split(resourceRaw, "=")
			item := Item{}
			if len(resource) == 2 {
				item.Resource = strings.TrimSpace(resource[1])
				item.Alias = strings.TrimSpace(resource[0])
			} else {
				item.Resource = strings.TrimSpace(resource[0])
			}
			c.Items = append(c.Items, item)
		}
	}
	return nil
}

func parseConfig(c *Config) error {
	if c.LogLevel == "" {
		c.LogLevel = "info"
	}
	lvl, err := logrus.ParseLevel(c.LogLevel)
	if err != nil {
		return err
	}
	logrus.SetLevel(lvl)

	if len(c.Items) == 0 {
		return errors.New("empty items list")
	}

	for _, item := range c.Items {
		hostPort := strings.Split(item.Resource, ":")
		if len(hostPort) != 2 {
			return fmt.Errorf("incorrect item: %+v", item)
		}
		portInt, err := strconv.Atoi(hostPort[1])
		if err != nil {
			return fmt.Errorf("incorrent port in item: %+v", item)
		}
		if len(item.Alias) == 0 {
			item.Alias = item.Resource
		}
		item.Host = hostPort[0]
		item.Port = portInt
	}

	return nil
}
