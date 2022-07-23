package main

import (
	"context"
	"errors"
	"fmt"
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
	if err != nil {
		return err
	}

	if status.Status != domain.HealthCheckStatusPass {
		return errors.New("server error")
	}

	return nil
}

// Init db - recreate bucket.
func (db *DB) Init() error {
	ctx := context.Background()

	bucket, err := db.client.BucketsAPI().FindBucketByName(ctx, db.bucket)
	if err != nil {
		return err
	} else {
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

// MapCO2 value to 7 categories, with 1 being the best.
func MapCO2(value int) int {
	if value < 0 && value <= 30 {
		return 1
	} else if value > 30 && value <= 70 {
		return 2
	} else if value > 70 && value <= 150 {
		return 3
	} else if value > 150 && value <= 200 {
		return 4
	} else if value > 200 && value <= 400 {
		return 5
	} else if value > 400 && value <= 800 {
		return 6
	}

	return 7
}

// MapAir quality value to 6 categories, with 1 being the best.
func MapAir(value int) int {
	if value < 0 && value <= 50 {
		return 1
	} else if value > 50 && value <= 100 {
		return 2
	} else if value > 100 && value <= 150 {
		return 3
	} else if value > 150 && value <= 200 {
		return 4
	} else if value > 200 && value <= 300 {
		return 5
	}

	return 6
}

// Add new data to db.
func (db *DB) Add(data *Data) {
	writeAPI := db.client.WriteAPI(db.org, db.bucket)
	p := influxdb2.NewPointWithMeasurement(data.DataType)

	switch data.DataType {
	case carbonMonoxide:
		p.AddField("value", data.Value)
		p.AddField("category", MapCO2(int(data.Value)))
	case airQuality:
		p.AddField("value", data.Value)
		p.AddField("category", MapAir(int(data.Value)))
	case raindrops:
		p.AddField("category", data.Value)
	case soilMoisture:
		p.AddField("category", data.Value)
	}

	writeAPI.WritePoint(p)
	writeAPI.Flush()
}

// Latest returns the latest data for the requested dataType.
func (db *DB) Latest(dataType string) (*DataResponse, error) {
	queryAPI := db.client.QueryAPI(db.org)
	query := fmt.Sprintf(`from(bucket:"%s")
			|> range(start: 0)
			|> filter(fn: (r) => r._measurement == "%s")
			|> last()`, db.bucket, dataType)

	dr := DataResponse{}

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil || result.Err() != nil {
		return &dr, err
	}

	for result.Next() {
		if result.Record().Field() == "category" {
			if res, ok := result.Record().Value().(int64); ok {
				dr.Category = int(res)
			}
		}

		if result.Record().Field() == "value" {
			if res, ok := result.Record().Value().(float64); ok {
				dr.Value = float32(res)
			}

			dr.DataType = result.Record().Measurement()
		}
	}

	return &dr, nil
}

// Last24H returns data for the last 24 hours for the requested dataType.
func (db *DB) Last24H(dataType string) ([]DataResponse, error) {
	queryAPI := db.client.QueryAPI(db.org)
	query := fmt.Sprintf(`from(bucket:"%s")
			|> range(start: -1d)
			|> filter(fn: (r) => r._measurement == "%s")`, db.bucket, dataType)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil || result.Err() != nil {
		return nil, err
	}

	var categories []int
	var values []float32

	for result.Next() {
		if result.Record().Field() == "category" {
			if res, ok := result.Record().Value().(int64); ok {
				categories = append(categories, int(res))
			} else {
				categories = append(categories, -1)
			}
		}

		if result.Record().Field() == "value" {
			if res, ok := result.Record().Value().(float64); ok {
				values = append(values, float32(res))
			} else {
				values = append(values, -1)
			}
		}
	}

	if len(categories) != len(values) {
		return nil, errors.New("length of categories is different than length of values")
	}

	var data []DataResponse

	for k, v := range categories {
		data = append(data, DataResponse{
			Data: Data{
				DataType: dataType,
				Value:    values[k],
			},
			Category: v,
		})
	}

	return data, nil
}
