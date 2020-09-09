package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/sirupsen/logrus"
)

const (
	defaultsListenAddr        = ":9407"
	defaultsMetricsPath       = "/metrics"
	defaultsLogLevel          = "info"
	defaultsLogFormat         = "text"
	defaultsConnectionTimeout = 500 * time.Millisecond
)

var (
	connectionTimeout   = flag.Duration("timeout", 0, "Connection timeout")
	logLevel            = flag.String("log-level", "", "Logging level")
	logFormat           = flag.String("log-format", "", "Logs format (text or json)")
	listenAddress       = flag.String("web.listen-address", "", "Listen address")
	metricsPath         = flag.String("web.telemetry-path", "", "Metrics path")
	resourcesCollection = flag.String("resources", "", "Resources list")
	configFile          = flag.String("config-file", "", "Configuration file in YAML format")
)

type StringMap map[string][]string

type Config struct {
	ConnectionTimeout time.Duration `yaml:"connectionTimeout"`
	LogLevel          string        `yaml:"logLevel"`
	LogFormat         string        `yaml:"logFormat"`
	ListenAddr        string        `yaml:"listenAddr"`
	MetricsPath       string        `yaml:"metricsPath"`
	RawItems          StringMap     `yaml:"resources"`
	Items             []Item        `yaml:"-"`
	File              string        `yaml:"-"`
}

func (i *StringMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var (
		value map[string][]string
		slice []string
	)
	// try to parse map[string][]string
	if err := unmarshal(&value); err != nil {
		// try to parse []string
		if err := unmarshal(&slice); err != nil {
			return err
		}
		value = map[string][]string{}
		value["all"] = slice
	}
	*i = (StringMap)(value)
	return nil
}

func LoadConfig() (*Config, error) {
	nc := &Config{
		File: *configFile,
	}
	if err := nc.LoadFromFile(); err != nil {
		return nil, err
	}
	if err := nc.LoadFromFlags(); err != nil {
		return nil, err
	}
	nc.SetEmptyToDefaults()
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
	if len(c.LogLevel) == 0 {
		c.LogLevel = defaultsLogLevel
	}
	if len(c.LogFormat) == 0 {
		c.LogFormat = defaultsLogFormat
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
	if len(*logFormat) != 0 {
		c.LogFormat = *logFormat
	}
	if len(*listenAddress) != 0 {
		c.ListenAddr = *listenAddress
	}
	if len(*metricsPath) != 0 {
		c.MetricsPath = *metricsPath
	}
	if len(*resourcesCollection) != 0 {
		resources := strings.Split(*resourcesCollection, ",")
		c.RawItems = make(map[string][]string)
		c.RawItems["all"] = resources
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
	lvl, err := logrus.ParseLevel(c.LogLevel)
	if err != nil {
		return err
	}
	logrus.SetLevel(lvl)

	switch c.LogFormat {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
		logrus.SetFormatter(&logrus.TextFormatter{})
	}

	if len(c.RawItems) == 0 {
		return errors.New("empty items list")
	}
	for group, items := range c.RawItems {
		for _, resource := range items {
			item, err := ParseResource(resource)
			if err != nil {
				return err
			}
			item.Group = group
			c.Items = append(c.Items, *item)
		}
	}
	return nil
}
