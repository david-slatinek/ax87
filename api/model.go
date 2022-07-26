package main

import (
	"fmt"
	"time"
)

const (
	// Carbon monoxide constant.
	carbonMonoxide = "carbonMonoxide"
	// Air quality constant.
	airQuality = "airQuality"
	// Raindrops constant.
	raindrops = "raindrops"
	// Soil moisture constant.
	soilMoisture = "soilMoisture"
)

// Data is a struct used for serializing data to the database and rarely from the database - only used with functions DB.Median, DB.Max, DB.Min.
type Data struct {
	// Presents data type. Use constants: carbonMonoxide, airQuality, raindrops, soilMoisture.
	DataType string
	// Sensor value.
	Value float32
	// When was the Value taken.
	TimeStamp time.Time
}

// DataResponse is a struct for serializing data from the database.
type DataResponse struct {
	// Embedded Data struct.
	Data
	// Data.Value category. Check functions MapCO2, MapAir, and MapValue.
	Category int
}

// String returns Data fields in a string.
func (d Data) String() string {
	return fmt.Sprintf("DataType: %s, Value: %f, Time: %v", d.DataType, d.Value, d.TimeStamp)
}

// String returns DataResponse fields in a string.
func (dr DataResponse) String() string {
	return fmt.Sprintf("%s, Category: %d", dr.Data, dr.Category)
}
