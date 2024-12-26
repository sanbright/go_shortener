package main

import (
	"log"
	"sanbright/go_shortener/internal/config"
	"testing"
)

func Test_initServer(t *testing.T) {
	configuration, err := config.NewConfig("localhost:80", "localhost", "", "", false, "", "", "")
	logger := setupLogger()
	if err != nil {
		log.Fatalf("Fatal configuration error: %s", err.Error())
	}

	tests := []struct {
		name   string
		config *config.Config
	}{
		{
			name:   "success",
			config: configuration,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := initServer(tt.config, logger)

			if result == nil {
				log.Fatalf("Error init application")
			}
		})
	}
}
