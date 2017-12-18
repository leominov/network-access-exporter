package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/sirupsen/logrus"
)

const (
	namespace    = "network_access"
	exporterName = "network_access_exporter"
)

var (
	showVersion    = flag.Bool("version", false, "Prints version information and exit")
	configPathFlag = flag.String("config-path", "/opt/prometheus/network-access-exporter.yaml", "Configuration file path")

	allowedResource = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "allowed"),
		"Was the last check successful",
		[]string{"resource", "alias"}, nil,
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

	logrus.Infof("Starting %s...", exporterName)

	cfg, err := LoadConfig(*configPathFlag)
	if err != nil {
		logrus.Fatal(err)
	}

	exporter := NewExporter(cfg)
	if err := prometheus.Register(exporter); err != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("Listen address: %s", cfg.ListenAddr)

	http.Handle(cfg.MetricPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>` + exporterName + ` v` + version.Version + `</title></head>
			<body>
			<h1>` + exporterName + ` v` + version.Version + `</h1>
			<p><a href='` + cfg.MetricPath + `'>Metrics</a></p>
			</body>
			</html>
		`))
	})
	logrus.Fatal(http.ListenAndServe(cfg.ListenAddr, nil))
}
