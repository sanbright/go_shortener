package config

import (
	"fmt"
	"net"
	"strings"
)

type Config struct {
	DomainAndPort DomainAndPort
	BaseURL       ExternalURL
}

type DomainAndPort struct {
	Domain string
	Port   string
}

type ExternalURL struct {
	URL string
}

func NewConfig(serverAddress string, baseURL string) *Config {
	if len(serverAddress) == 0 {
		serverAddress = "localhost:8080"
	}

	if len(baseURL) == 0 {
		baseURL = "http://localhost:8080"
	}

	var domainAndPort DomainAndPort
	var externalURL ExternalURL

	_ = externalURL.Set(baseURL)
	err := domainAndPort.Set(serverAddress)

	if err != nil {
		return nil
	}

	return &Config{
		DomainAndPort: domainAndPort,
		BaseURL:       externalURL,
	}
}

func (dap *DomainAndPort) String() string {
	arr := make([]string, 0)
	arr = append(arr, dap.Domain, dap.Port)

	return fmt.Sprint(strings.Join(arr, ":"))
}

func (eu *ExternalURL) String() string {
	return eu.URL
}

func (eu *ExternalURL) Set(value string) error {
	eu.URL = value

	return nil
}

func (dap *DomainAndPort) Set(value string) error {
	domain, port, err := net.SplitHostPort(value)
	if err != nil {
		return err
	}

	dap.Domain = domain
	dap.Port = port

	return nil
}
