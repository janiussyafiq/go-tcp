package main

import (
	"flag"
	"fmt"
	"strings"
)

type Config struct {
	ProxyPort string
	Backends  []string
}

func ParseConfig() (*Config, error) {
	proxyPort := flag.String("port", "8080", "Port for the proxy server")
	backends := flag.String("backends", "http://localhost:8081", "Comma-separated list of backend URLs")

	flag.Parse()

	// Split the backends by comma
	backendList := strings.Split(*backends, ",")

	// Trim spaces
	for i, backend := range backendList {
		backendList[i] = strings.TrimSpace(backend)
	}

	if len(backendList) == 0 {
		return nil, fmt.Errorf("at least one backend is required")
	}

	return &Config{
		ProxyPort: *proxyPort,
		Backends:  backendList,
	}, nil
}
