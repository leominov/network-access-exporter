package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

const (
	defaultsListenAddr        = ":9407"
	defaultsMetricPath        = "/metrics"
	defaultsConnectionTimeout = 500 * time.Millisecond
)

type Config struct {
	ConnectionTimeout time.Duration `yaml:"connectionTimeout"`
	LogLevel          string        `yaml:"logLevel"`
	ListenAddr        string        `yaml:"listenAddr"`
	MetricPath        string        `yaml:"metricPath"`
	Items             []*Item       `yaml:"items"`
}

type Item struct {
	Alias    string `yaml:"alias"`
	Resource string `yaml:"resource"`
	Host     string `yaml:"-"`
	Port     int    `yaml:"-"`
}

func LoadConfig(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %v", path, err)
	}
	nc := &Config{}
	if err := yaml.Unmarshal(b, nc); err != nil {
		return nil, fmt.Errorf("error unmarshaling %s: %v", path, err)
	}
	if err := parseConfig(nc); err != nil {
		return nil, err
	}
	return nc, nil
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

	if len(c.ListenAddr) == 0 {
		c.ListenAddr = defaultsListenAddr
	}

	if len(c.MetricPath) == 0 {
		c.MetricPath = defaultsMetricPath
	}

	if len(c.Items) == 0 {
		return errors.New("Empty items list")
	}

	if c.ConnectionTimeout.Nanoseconds() == 0 {
		c.ConnectionTimeout = defaultsConnectionTimeout
	}

	for _, item := range c.Items {
		hostPort := strings.Split(item.Resource, ":")
		if len(hostPort) != 2 {
			return fmt.Errorf("Incorrect item: %s", item)
		}
		portInt, err := strconv.Atoi(hostPort[1])
		if err != nil {
			return fmt.Errorf("Incorrent port in item: %s", item)
		}
		item.Host = hostPort[0]
		item.Port = portInt
	}

	return nil
}
