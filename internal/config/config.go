package config

import (
	"fmt"
	"net"
	"strings"
)

type Config struct {
	DomainAndPort DomainAndPort
	BaseUrl       ExternalURL
}

type DomainAndPort struct {
	Domain string
	Port   string
}

type ExternalURL struct {
	URL string
}

func NewConfig(serverAddress string, baseUrl string) *Config {
	if len(serverAddress) == 0 {
		serverAddress = "localhost:8080"
	}

	if len(baseUrl) == 0 {
		baseUrl = "http://localhost:8080"
	}

	var domainAndPort DomainAndPort
	var externalUrl ExternalURL

	_ = externalUrl.Set(baseUrl)
	err := domainAndPort.Set(serverAddress)

	if err != nil {
		return nil
	}

	return &Config{
		DomainAndPort: domainAndPort,
		BaseUrl:       externalUrl,
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
