package main

import (
	"context"
	"errors"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	"github.com/joho/godotenv"
	"os"
)

// loadENV Loads environment variables from '.env' file.
func loadENV() error {
	return godotenv.Load(".env")
}

// connectToDB Connects to the influxdb server.
func connectToDB() (influxdb2.Client, error) {
	token := os.Getenv("INFLUXDB_TOKEN")
	if token == "" {
		return nil, errors.New("empty INFLUXDB_TOKEN")
	}

	url := os.Getenv("INFLUXDB_URL")
	if url == "" {
		return nil, errors.New("empty INFLUXDB_URL")
	}

	client := influxdb2.NewClientWithOptions(url, token, influxdb2.DefaultOptions().SetUseGZip(true))

	status, err := client.Health(context.Background())
	if status.Status != domain.HealthCheckStatusPass {
		return nil, err
	}

	return client, nil
}
