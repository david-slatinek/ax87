package main

import (
	"context"
	"errors"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	"os"
	"time"
)

// DB is a struct for database management.
type DB struct {
	// influxdb2.Client for influxdb communication.
	client influxdb2.Client
	// Database connection string.
	url string
	// Database token.
	token string
	// Database organization name.
	org string
	// Database bucket name.
	bucket string
}

// LoadFields loads DB fields - url, token, org, bucket - from the environment.
func (db *DB) LoadFields() {
	db.url = os.Getenv("INFLUXDB_URL")
	db.token = os.Getenv("INFLUXDB_TOKEN")
	db.org = os.Getenv("ORGANIZATION")
	db.bucket = os.Getenv("BUCKET")
}

// Connect to the influxdb server.
func (db *DB) Connect() error {
	if db == nil {
		return errors.New("db can't be nil")
	}

	db.client = influxdb2.NewClientWithOptions(db.url, db.token,
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

// Init db - recreate (delete) bucket.
func (db *DB) Init() error {
	if db == nil {
		return errors.New("db can't be nil")
	}

	if db.client == nil {
		return errors.New("db.client can't be nil")
	}

	ctx := context.Background()

	bucket, err := db.client.BucketsAPI().FindBucketByName(ctx, db.bucket)
	if err != nil {
		return err
	}

	if err := db.client.BucketsAPI().DeleteBucketWithID(context.Background(), *bucket.Id); err != nil {
		return err
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
	if value <= 30 {
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

// MapValue from one range to another.
func MapValue(x int, inMin int, inMax int, outMin int, outMax int) int {
	return (x-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}

// Add new data to db.
func (db *DB) Add(data *Data) {
	if db == nil || data == nil || db.client == nil {
		return
	}
	writeAPI := db.client.WriteAPI(db.org, db.bucket)
	p := influxdb2.NewPointWithMeasurement(data.DataType).AddField("value", data.Value).SetTime(data.TimeStamp)

	switch data.DataType {
	case carbonMonoxide:
		p.AddField("category", MapCO2(int(data.Value)))
	case airQuality:
		p.AddField("category", MapAir(int(data.Value)))
	case raindrops:
		if int(data.Value) < 0 || int(data.Value) > 1024 {
			return
		}
		p.AddField("category", MapValue(int(data.Value), 0, 1024, 1, 4))
	case soilMoisture:
		if int(data.Value) < 238 || int(data.Value) > 489 {
			return
		}
		p.AddField("category", MapValue(int(data.Value), 489, 238, 0, 100))
	default:
		return
	}

	writeAPI.WritePoint(p)
	writeAPI.Flush()
}

// Latest returns the latest data for the requested dataType.
func (db *DB) Latest(dataType string) (*DataResponse, error) {
	if db == nil {
		return nil, errors.New("db can't be nil")
	}

	if db.client == nil {
		return nil, errors.New("db.client can't be nil")
	}

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

	isData := false

	for result.Next() {
		isData = true

		if result.Record().Field() == "category" {
			if res, ok := result.Record().Value().(int64); ok {
				dr.Category = int(res)
			} else {
				dr.Category = -1
			}
		}

		if result.Record().Field() == "value" {
			if res, ok := result.Record().Value().(float64); ok {
				dr.Value = float32(res)
			} else {
				dr.Value = -1
			}
		}
	}

	if !isData {
		return &dr, errors.New("no data")
	}

	dr.DataType = result.Record().Measurement()
	dr.TimeStamp = result.Record().Time()

	return &dr, nil
}

// Last24H returns data for the last 24 hours for the requested dataType.
func (db *DB) Last24H(dataType string) (*[]DataResponse, error) {
	if db == nil {
		return nil, errors.New("db can't be nil")
	}

	if db.client == nil {
		return nil, errors.New("db.client can't be nil")
	}

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
	var tm []time.Time

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
			tm = append(tm, result.Record().Time())
		}
	}

	if len(categories) != len(values) || len(categories) != len(tm) {
		return nil, errors.New("different lengths of categories, values, and timestamps")
	}

	var data []DataResponse

	for index, value := range categories {
		data = append(data, DataResponse{
			Data: Data{
				DataType:  dataType,
				Value:     values[index],
				TimeStamp: tm[index],
			},
			Category: value,
		})
	}

	return &data, nil
}

// RetrieveData returns data for the given query. Only used for Median, Max and Min to prevent code duplication.
func (db *DB) RetrieveData(query string) (*Data, error) {
	if db == nil {
		return nil, errors.New("db can't be nil")
	}

	if db.client == nil {
		return nil, errors.New("db.client can't be nil")
	}

	queryAPI := db.client.QueryAPI(db.org)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil || result.Err() != nil {
		return nil, err
	}

	dr := Data{}
	isData := false

	for result.Next() {
		isData = true
		if res, ok := result.Record().Value().(float64); ok {
			dr.Value = float32(res)
		} else {
			dr.Value = -1
		}
	}

	if !isData {
		return &dr, errors.New("no data")
	}

	dr.DataType = result.Record().Measurement()
	dr.TimeStamp = result.Record().Time()

	return &dr, nil
}

// Median returns median data for the last 24 hours for the requested dataType.
func (db *DB) Median(dataType string) (*Data, error) {
	if db == nil {
		return nil, errors.New("db can't be nil")
	}

	if db.client == nil {
		return nil, errors.New("db.client can't be nil")
	}

	return db.RetrieveData(fmt.Sprintf(`from(bucket:"%s")
			|> range(start: -1d)
			|> filter(fn: (r) => r._measurement == "%s" and r._field == "value")
			|> median(method: "exact_selector")`, db.bucket, dataType))
}

// Max returns maximum data for the last 24 hours for the requested dataType.
func (db *DB) Max(dataType string) (*Data, error) {
	if db == nil {
		return nil, errors.New("db can't be nil")
	}

	if db.client == nil {
		return nil, errors.New("db.client can't be nil")
	}

	return db.RetrieveData(fmt.Sprintf(`from(bucket:"%s")
			|> range(start: -1d)
			|> filter(fn: (r) => r._measurement == "%s" and r._field == "value")
			|> max()`, db.bucket, dataType))
}

// Min returns minimum data for the last 24 hours for the requested dataType.
func (db *DB) Min(dataType string) (*Data, error) {
	if db == nil {
		return nil, errors.New("db can't be nil")
	}

	if db.client == nil {
		return nil, errors.New("db.client can't be nil")
	}

	return db.RetrieveData(fmt.Sprintf(`from(bucket:"%s")
			|> range(start: -1d)
			|> filter(fn: (r) => r._measurement == "%s" and r._field == "value")
			|> min()`, db.bucket, dataType))
}
