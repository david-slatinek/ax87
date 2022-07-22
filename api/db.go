package main

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	"os"
)

type DB struct {
	client influxdb2.Client
}

// connectToDB Connects to the influxdb server.
func (db *DB) connectToDB() error {
	db.client = influxdb2.NewClientWithOptions(os.Getenv("INFLUXDB_URL"), os.Getenv("INFLUXDB_TOKEN"),
		influxdb2.DefaultOptions().SetUseGZip(true))

	status, err := db.client.Health(context.Background())
	if status.Status != domain.HealthCheckStatusPass {
		return err
	}

	return nil
}
