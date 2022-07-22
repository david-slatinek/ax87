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
		if err := db.client.BucketsAPI().DeleteBucketWithID(context.Background(), *bucket.Id); err != nil {
			return err
		}
	}

	org, err := db.client.OrganizationsAPI().FindOrganizationByName(ctx, db.org)
	if err != nil {
		return err
	}

	_, err = db.client.BucketsAPI().CreateBucketWithNameWithID(ctx, *org.Id, db.bucket)
	return err
}

// MapValue from one range to another.
func MapValue(x int, inMin int, inMax int, outMin int, outMax int) int {
	return (x-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}

// Add new data to db.
func (db *DB) Add(data *Data) {
	writeAPI := db.client.WriteAPI(db.org, db.bucket)
	p := influxdb2.NewPointWithMeasurement(data.DataType.String()).AddField("value", data.Value)

	if data.DataType == carbonMonoxide || data.DataType == airQuality {
		var category int

		switch data.DataType {
		case carbonMonoxide:
			category = MapValue(int(data.Value), 0, 13000, 0, 3)
		case airQuality:
			category = MapValue(int(data.Value), 0, 500, 0, 5)
		}
		p.AddField("category", category)
	}

	writeAPI.WritePoint(p)
	writeAPI.Flush()
}
