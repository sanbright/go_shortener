// Package config пакет конфигурации
package config

import (
	"cmp"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
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
	DomainAndPort DomainAndPort `json:"server_address"`
	// BaseURL - Базовый URL необходим для формирования коротких ссылок
	BaseURL ExternalURL `json:"base_url"`
	// StoragePath - путь до файла, кторый используется как файловое хранилище
	StoragePath string `json:"file_storage_path"`
	// DatabaseDSN - DSN для использования подключения к СУБД
	DatabaseDSN string `json:"database_dsn"`
	// HTTPS - Использование https
	HTTPS bool `json:"enable_https"`
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
func NewConfig(serverAddress string, baseURL string, storagePath string, databaseDSN string, HTTPS bool, configFile string) (*Config, error) {

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

	flag.Var(&domainAndPort, "a", "listen host and port")
	flag.Var(&externalURL, "b", "domain in short link")
	flag.StringVar(&storagePathConf, "f", "", "file storage path")
	flag.StringVar(&databaseDSNConf, "d", "", "database storage")
	flag.BoolVar(&HTTPSConf, "s", HTTPS, "HTTPS Enable")
	flag.Parse()

	fileConfig := readConfig(configFile)

	fmt.Printf("Config: %+v\n", fileConfig)

	var defC = &DomainAndPort{}
	err := defC.Set(serverAddress)
	if err != nil {
		return nil, err
	}

	var defEx = &ExternalURL{}
	err = defEx.Set(baseURL)
	if err != nil {
		return nil, err
	}

	return &Config{
		DomainAndPort: cmp.Or(domainAndPort, fileConfig.DomainAndPort, *defC),
		BaseURL:       cmp.Or(externalURL, fileConfig.BaseURL, *defEx),
		StoragePath:   cmp.Or(storagePathConf, fileConfig.StoragePath, storagePath),
		DatabaseDSN:   cmp.Or(databaseDSNConf, fileConfig.DatabaseDSN, databaseDSN),
		HTTPS:         cmp.Or(HTTPSConf, fileConfig.HTTPS, false),
	}, nil
}

// readConfig - загрузка конфигурации из файла
func readConfig(configFile string) *Config {
	fileConfig := &Config{}
	if configFile != "" {
		rawContent, err := os.ReadFile(configFile)
		if err != nil {
			log.Fatal(err)
		}

		if err = json.Unmarshal(rawContent, fileConfig); err != nil {
			log.Fatal(err)
		}
	}

	return fileConfig
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

// UnmarshalJSON - Unmarshal конфига
func (dap *DomainAndPort) UnmarshalJSON(data []byte) error {

	domainAndPort, err := strconv.Unquote(string(data))

	if err != nil {
		return err
	}

	dap.Set(domainAndPort)
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

// UnmarshalJSON - Unmarshal конфига
func (eu *ExternalURL) UnmarshalJSON(data []byte) error {
	externalURL, err := strconv.Unquote(string(data))

	if err != nil {
		return err
	}

	eu.Set(externalURL)

	return nil
}
