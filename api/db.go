package main

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	"os"
)

// DB struct for db management.
type DB struct {
	client influxdb2.Client
	org    string
	bucket string
}

// Connect to the influxdb server.
func (db *DB) Connect() error {
	db.client = influxdb2.NewClientWithOptions(os.Getenv("INFLUXDB_URL"), os.Getenv("INFLUXDB_TOKEN"),
		influxdb2.DefaultOptions().SetUseGZip(true))

	status, err := db.client.Health(context.Background())
	if status.Status != domain.HealthCheckStatusPass {
		return err
	}

	return nil
}

// Init db - recreate bucket.
func (db *DB) Init() error {
	ctx := context.Background()

	bucket, err := db.client.BucketsAPI().FindBucketByName(ctx, db.bucket)
	if err == nil {
		err := db.client.BucketsAPI().DeleteBucketWithID(context.Background(), *bucket.Id)
		if err != nil {
			return err
		}
	}

	org, err := db.client.OrganizationsAPI().FindOrganizationByName(ctx, db.org)
	if err != nil {
		return err
	}

	_, err = db.client.BucketsAPI().CreateBucketWithNameWithID(ctx, *org.Id, db.bucket)

	if err != nil {
		return err
	}

	return nil
}

// Add new data to db.
func (db *DB) Add(data *Data) {
	writeAPI := db.client.WriteAPI(db.org, db.bucket)
	p := influxdb2.NewPointWithMeasurement("measure").
		AddTag("type", data.DataType.String()).AddField("value", data.Value)

	writeAPI.WritePoint(p)
	writeAPI.Flush()
}
