// Package config пакет конфигурации
package config

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// Дефолтная конфигурация приложения
const (
	DefaultServerAddress string = "localhost:8080"
	DefaultBaseURL       string = "http://localhost:8080"
	DefaultStoragePath   string = "/test.db"
	DefaultDatabaseDSN   string = ""
)

// Config структура конфигурации
type Config struct {
	// DomainAndPort - домен и порт для запуска web сервера
	DomainAndPort DomainAndPort
	// BaseURL - Базовый URL необходим для формирования коротких ссылок
	BaseURL ExternalURL
	// StoragePath - путь до файла, кторый используется как файловое хранилище
	StoragePath string
	// DatabaseDSN - DSN для использования подключения к СУБД
	DatabaseDSN string
	// HTTPS - Использование https
	HTTPS bool
}

// DomainAndPort - Домен и порт
type DomainAndPort struct {
	// Domain - домен
	Domain string
	// Port - порт
	Port string
}

// ExternalURL - Внешний УРЛ
type ExternalURL struct {
	URL string
}

// NewConfig Конструктор инициализация конфиругации
func NewConfig(serverAddress string, baseURL string, storagePath string, databaseDSN string, HTTPS bool) (*Config, error) {

	if len(serverAddress) == 0 {
		serverAddress = DefaultServerAddress
	}

	if len(baseURL) == 0 {
		baseURL = DefaultBaseURL
	}

	if len(databaseDSN) == 0 {
		databaseDSN = DefaultDatabaseDSN
	}

	if len(storagePath) == 0 {
		basePath, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}
		storagePath = basePath + DefaultStoragePath
	}

	var domainAndPort DomainAndPort
	var externalURL ExternalURL
	var storagePathConf string
	var databaseDSNConf string
	var HTTPSConf bool

	err := externalURL.Set(baseURL)
	if err != nil {
		return nil, err
	}

	err = domainAndPort.Set(serverAddress)
	if err != nil {
		return nil, err
	}

	fmt.Printf("HTTPSConf value \n%s", HTTPSConf)

	flag.Var(&domainAndPort, "a", "listen host and port")
	flag.Var(&externalURL, "b", "domain in short link")
	flag.StringVar(&storagePathConf, "f", storagePath, "file storage path")
	flag.StringVar(&databaseDSNConf, "d", databaseDSN, "database storage")
	flag.BoolVar(&HTTPSConf, "s", HTTPS, "HTTPS Enable")
	flag.Parse()

	fmt.Printf("HTTPSConf post value \n%s", HTTPSConf)

	return &Config{
		DomainAndPort: domainAndPort,
		BaseURL:       externalURL,
		StoragePath:   storagePathConf,
		DatabaseDSN:   databaseDSNConf,
		HTTPS:         HTTPSConf,
	}, nil
}

// String - преобразование домена и порта в строку
func (dap *DomainAndPort) String() string {
	arr := make([]string, 0)
	arr = append(arr, dap.Domain, dap.Port)

	return fmt.Sprint(strings.Join(arr, ":"))
}

// Set - преобразование строки в домен и порт
func (dap *DomainAndPort) Set(value string) error {
	domain, port, err := net.SplitHostPort(value)
	if err != nil {
		return err
	}

	dap.Domain = domain
	dap.Port = port

	return nil
}

// String - преобразование значения URL в строку
func (eu *ExternalURL) String() string {
	return eu.URL
}

// Set - установка значение URL
func (eu *ExternalURL) Set(value string) error {
	eu.URL = value

	return nil
}
