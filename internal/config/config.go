package config

import (
	"flag"
	"fmt"
	"net"
	"strings"
)

const (
	DefaultServerAddress string = "localhost:8080"
	DefaultBaseURL       string = "http://localhost:8080"
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

func NewConfig(serverAddress string, baseURL string) (*Config, error) {

	if len(serverAddress) == 0 {
		serverAddress = DefaultServerAddress
	}

	if len(baseURL) == 0 {
		baseURL = DefaultBaseURL
	}

	var domainAndPort DomainAndPort
	var externalURL ExternalURL

	err := externalURL.Set(baseURL)
	if err != nil {
		return nil, err
	}

	err = domainAndPort.Set(serverAddress)
	if err != nil {
		return nil, err
	}

	flag.Var(&domainAndPort, "a", "listen host and port")
	flag.Var(&externalURL, "b", "domain in short link")
	flag.Parse()

	return &Config{
		DomainAndPort: domainAndPort,
		BaseURL:       externalURL,
	}, nil
}

func (dap *DomainAndPort) String() string {
	arr := make([]string, 0)
	arr = append(arr, dap.Domain, dap.Port)

	return fmt.Sprint(strings.Join(arr, ":"))
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

func (eu *ExternalURL) String() string {
	return eu.URL
}

func (eu *ExternalURL) Set(value string) error {
	eu.URL = value

	return nil
}
