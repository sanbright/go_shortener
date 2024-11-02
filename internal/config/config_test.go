package config

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
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
		want          want
	}{
		{
			name:          "Standard_configuration",
			serverAddress: "localhost:8081",
			baseURL:       "http://example.ru",
			storagePath:   "/test.db",
			databaseDSN:   "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=go_mark sslmode=disable",
			want: want{
				databaseDSN:       "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=go_mark sslmode=disable",
				storagePath:       "/test.db",
				serverAddressHost: "localhost",
				serverAddressPort: "8081",
				baseURL:           "http://example.ru",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := NewConfig(tt.serverAddress, tt.baseURL, tt.storagePath, tt.databaseDSN)

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
