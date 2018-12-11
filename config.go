package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

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
	configFile          = flag.String("config-file", "", "Configuration file in YAML format")
)

type Config struct {
	ConnectionTimeout time.Duration `yaml:"connectionTimeout"`
	LogLevel          string        `yaml:"logLevel"`
	ListenAddr        string        `yaml:"listenAddr"`
	MetricsPath       string        `yaml:"metricsPath"`
	Items             []Item        `yaml:"items"`
}

func LoadConfig() (*Config, error) {
	nc := &Config{}
	if err := nc.LoadFromFile(); err != nil {
		return nil, err
	}
	if err := nc.LoadFromFlags(); err != nil {
		return nil, err
	}
	if err := parseConfig(nc); err != nil {
		return nil, err
	}
	nc.SetEmptyToDefaults()
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
			item := Item{
				Resource: resourceRaw,
			}
			c.Items = append(c.Items, item)
		}
	}
	return nil
}

func (c *Config) LoadFromFile() error {
	if len(*configFile) == 0 {
		return nil
	}
	b, err := ioutil.ReadFile(*configFile)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(b, &c); err != nil {
		return err
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
	for id, item := range c.Items {
		hostPort := strings.Split(item.Resource, ":")
		if len(hostPort) != 2 {
			return fmt.Errorf("incorrect item: %+v", item)
		}
		portInt, err := strconv.Atoi(hostPort[1])
		if err != nil {
			return fmt.Errorf("incorrent port in item: %+v", item)
		}
		item.Host = hostPort[0]
		item.Port = portInt
		c.Items[id] = item
	}
	return nil
}
