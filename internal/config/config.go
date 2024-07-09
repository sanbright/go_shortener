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

func NewConfig() *Config {
	return &Config{
		DomainAndPort: DomainAndPort{Domain: "localhost", Port: "8080"},
		BaseUrl:       ExternalURL{URL: "http://localhost:8080"},
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
