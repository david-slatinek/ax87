package db

import (
	"api/model"
	"api/util"
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

// Load DB fields - url, token, org, bucket - from the environment.
func (db *DB) Load() {
	db.url = os.Getenv("INFLUXDB_URL")
	db.token = os.Getenv("INFLUXDB_TOKEN")
	db.org = os.Getenv("ORGANIZATION")
	db.bucket = os.Getenv("BUCKET")
}

// Close db connection.
func (db *DB) Close() {
	if db.client != nil {
		db.client.Close()
	}
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

// Add new data to the db.
func (db *DB) Add(data *model.Data) {
	if db == nil || data == nil || db.client == nil {
		return
	}
	writeAPI := db.client.WriteAPI(db.org, db.bucket)
	p := influxdb2.NewPointWithMeasurement(data.DataType).AddField("value", data.Value).SetTime(data.Timestamp.Round(0))

	switch data.DataType {
	case util.Raindrops:
		if int(data.Value) < 0 || int(data.Value) > 1024 {
			return
		}
	case util.SoilMoisture:
		if int(data.Value) < 238 || int(data.Value) > 489 {
			return
		}
	}

	p.AddField("category", util.GetCategory(int(data.Value), data.DataType))

	writeAPI.WritePoint(p)
	writeAPI.Flush()
}

// Latest returns the latest data for the requested dataType.
func (db *DB) Latest(dataType string) (*model.DataResponse, error) {
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

	dr := model.DataResponse{}

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
	dr.Timestamp = result.Record().Time().Local()

	return &dr, nil
}

// Last24H returns data for the last 24 hours for the requested dataType.
func (db *DB) Last24H(dataType string) (*[]model.DataResponse, error) {
	if db == nil {
		return nil, errors.New("db can't be nil")
	}

	if db.client == nil {
		return nil, errors.New("db.client can't be nil")
	}

	queryAPI := db.client.QueryAPI(db.org)
	query := fmt.Sprintf(`from(bucket:"%s")
			|> range(start: -1d)
			|> filter(fn: (r) => r._measurement == "%s")
			|> sort(columns: ["_time"], desc: true)`, db.bucket, dataType)

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
			tm = append(tm, result.Record().Time().Local())
		}
	}

	if len(categories) != len(values) || len(categories) != len(tm) {
		return nil, errors.New("different lengths of categories, values, and timestamps")
	}

	var data []model.DataResponse

	for index, value := range categories {
		data = append(data, model.DataResponse{
			Data: model.Data{
				DataType:  dataType,
				Value:     values[index],
				Timestamp: tm[index],
			},
			Category: value,
		})
	}

	return &data, nil
}

// RetrieveData returns data for the given query. Only used for Median, Max and Min to prevent code duplication.
func (db *DB) RetrieveData(query, dataType string) (*model.DataResponse, error) {
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

	dr := model.DataResponse{}
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
	dr.Timestamp = result.Record().Time().Local()
	dr.Category = util.GetCategory(int(dr.Value), dataType)

	return &dr, nil
}

// Median returns median data for the last 24 hours for the requested dataType.
func (db *DB) Median(dataType string) (*model.DataResponse, error) {
	if db == nil {
		return nil, errors.New("db can't be nil")
	}

	if db.client == nil {
		return nil, errors.New("db.client can't be nil")
	}

	return db.RetrieveData(fmt.Sprintf(`from(bucket:"%s")
			|> range(start: -1d)
			|> filter(fn: (r) => r._measurement == "%s" and r._field == "value")
			|> median(method: "exact_selector")`, db.bucket, dataType), dataType)
}

// Max returns maximum data for the last 24 hours for the requested dataType.
func (db *DB) Max(dataType string) (*model.DataResponse, error) {
	if db == nil {
		return nil, errors.New("db can't be nil")
	}

	if db.client == nil {
		return nil, errors.New("db.client can't be nil")
	}

	return db.RetrieveData(fmt.Sprintf(`from(bucket:"%s")
			|> range(start: -1d)
			|> filter(fn: (r) => r._measurement == "%s" and r._field == "value")
			|> max()`, db.bucket, dataType), dataType)
}

// Min returns minimum data for the last 24 hours for the requested dataType.
func (db *DB) Min(dataType string) (*model.DataResponse, error) {
	if db == nil {
		return nil, errors.New("db can't be nil")
	}

	if db.client == nil {
		return nil, errors.New("db.client can't be nil")
	}

	return db.RetrieveData(fmt.Sprintf(`from(bucket:"%s")
			|> range(start: -1d)
			|> filter(fn: (r) => r._measurement == "%s" and r._field == "value")
			|> min()`, db.bucket, dataType), dataType)
}
