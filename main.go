package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/sirupsen/logrus"
)

const (
	ConfigPath = "/config"

	namespace    = "network_access"
	exporterName = "network_access_exporter"
)

var (
	showVersion = flag.Bool("version", false, "Prints version information and exit")
	lintConfig  = flag.Bool("lint", false, "Check configuration and exit")

	allowedResource = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "allowed"),
		"Was the last check successful",
		[]string{"resource", "meta_group", "ip", "ipv6"}, nil,
	)

	loopkupError = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "lookup_error"),
		"Error with getting IP address of resource",
		[]string{"resource", "meta_group"}, nil,
	)

	loopkupDuration = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "lookup_duration_seconds"),
		"Time spent for resource lookup in seconds",
		[]string{"resource", "meta_group"}, nil,
	)

	dialDuration = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "dial_duration_seconds"),
		"Time spent for TCP dial in seconds",
		[]string{"resource", "meta_group", "ip", "ipv6"}, nil,
	)
)

func init() {
	prometheus.MustRegister(version.NewCollector(exporterName))
}

func versionInfo() {
	fmt.Println(version.Print(exporterName))
	os.Exit(0)
}

func main() {
	flag.Parse()

	if *showVersion == true {
		versionInfo()
	}

	if !*lintConfig {
		logrus.Infof("Starting %s %s...", exporterName, version.Version)
	}

	cfg, err := LoadConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	if *lintConfig {
		logrus.Infof("Resources: %d", len(cfg.Items))
		logrus.Info("Look's good")
		return
	}

	exporter := NewExporter(cfg)
	if err := prometheus.Register(exporter); err != nil {
		logrus.Fatal(err)
	}

	if len(cfg.File) != 0 {
		logrus.Infof("Configuration file: %s", cfg.File)
	}
	logrus.Infof("Listen address: %s", cfg.ListenAddr)
	logrus.Infof("Connection timeout: %s", cfg.ConnectionTimeout.String())
	logrus.Debugf("Resources: %#v", cfg.Items)

	http.Handle(cfg.MetricsPath, promhttp.Handler())
	http.HandleFunc(ConfigPath, func(w http.ResponseWriter, r *http.Request) {
		b, err := yaml.Marshal(cfg)
		if err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
		w.Write([]byte(b))
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>` + exporterName + ` v` + version.Version + `</title></head>
			<body>
			<h1>` + exporterName + ` v` + version.Version + `</h1>
			<p><a href='` + cfg.MetricsPath + `'>Metrics</a></p>
			<p><a href='` + ConfigPath + `'>Configuration</a></p>
			</body>
			</html>
		`))
	})
	logrus.Fatal(http.ListenAndServe(cfg.ListenAddr, nil))
}
