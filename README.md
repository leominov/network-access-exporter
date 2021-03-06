# Just another Prometheus exporter

[![Build Status](https://travis-ci.com/leominov/network-access-exporter.svg?branch=master)](https://travis-ci.com/leominov/network-access-exporter)

## Usage

### Flags

```
Usage of ./network-access-exporter:
  -config-file string
    	Configuration file in YAML format
  -lint
    	Check configuration and exit
  -log-format string
        Logs format (text or json)
  -log-level string
    	Logging level
  -resources string
    	Resources list
  -timeout duration
    	Connection timeout
  -version
    	Prints version information and exit
  -web.listen-address string
    	Listen address
  -web.telemetry-path string
    	Metrics path
```

### Configuration file

See [example](config.example.yaml).

## Metrics

* `network_access_allowed` – Was the last check successful
* `network_access_lookup_error` – Lookup error by resource
* `network_access_lookup_duration_seconds` – Time spent for resource lookup in seconds
* `network_access_dial_duration_seconds` – Time spent for TCP dial in seconds

## Links

- https://hub.docker.com/repository/docker/leominov/network-access-exporter
- https://github.com/leominov/network-access-exporter/releases
