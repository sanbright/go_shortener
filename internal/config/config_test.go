package config

import (
	"flag"
	"log"
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {

	basePath, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	type want struct {
		databaseDSN       string
		storagePath       string
		serverAddressHost string
		serverAddressPort string
		baseURL           string
	}

	tests := []struct {
		name          string
		serverAddress string
		baseURL       string
		storagePath   string
		databaseDSN   string
		fileConfig    string
		want          want
	}{
		{
			name:          "Standard_configuration",
			serverAddress: "localhost:8081",
			baseURL:       "http://example.ru",
			storagePath:   "/test.db",
			databaseDSN:   "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=go_mark sslmode=disable",
			fileConfig:    "",
			want: want{
				databaseDSN:       "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=go_mark sslmode=disable",
				storagePath:       "/test.db",
				serverAddressHost: "localhost",
				serverAddressPort: "8081",
				baseURL:           "http://example.ru",
			},
		},
		{
			name:          "Default_configuration",
			serverAddress: "",
			baseURL:       "",
			storagePath:   "",
			databaseDSN:   "",
			fileConfig:    "",
			want: want{
				databaseDSN:       "",
				storagePath:       basePath + "/test.db",
				serverAddressHost: "localhost",
				serverAddressPort: "8080",
				baseURL:           "http://localhost:8080",
			},
		},
		{
			name:          "File_configuration",
			serverAddress: "",
			baseURL:       "",
			storagePath:   "",
			databaseDSN:   "",
			fileConfig:    basePath + "/../../config.json.dist",
			want: want{
				databaseDSN:       "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=go_mark sslmode=enable",
				storagePath:       basePath + "/test.db",
				serverAddressHost: "localhost",
				serverAddressPort: "8081",
				baseURL:           "http://localhost:8081",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet("-a localhost:8081 -f test.db", flag.ContinueOnError)

			config, err := NewConfig(tt.serverAddress, tt.baseURL, tt.storagePath, tt.databaseDSN, false, tt.fileConfig)

			if err != nil {
				t.Errorf("%v: ERROR '%v'", tt.name, err.Error())
			}

			if code := tt.want.databaseDSN; code != config.DatabaseDSN {
				t.Errorf("%v: DatabaseDSN = '%v', want = '%v'", tt.name, config.DatabaseDSN, tt.want.databaseDSN)
			}

			if code := tt.want.serverAddressHost; code != config.DomainAndPort.Domain {
				t.Errorf("%v: DomainAndPort.Domain = '%v', want = '%v'", tt.name, config.DomainAndPort.Domain, tt.want.serverAddressHost)
			}

			if code := tt.want.serverAddressPort; code != config.DomainAndPort.Port {
				t.Errorf("%v: DomainAndPort.Port = '%v', want = '%v'", tt.name, config.DomainAndPort.Port, tt.want.serverAddressPort)
			}

			if code := tt.want.storagePath; code != config.StoragePath {
				t.Errorf("%v: StoragePath.Port = '%v', want = '%v'", tt.name, config.StoragePath, tt.want.storagePath)
			}

			if code := tt.want.baseURL; code != config.BaseURL.String() {
				t.Errorf("%v: BaseURL = '%v', want = '%v'", tt.name, config.BaseURL.String(), tt.want.baseURL)
			}
		})
	}
}
