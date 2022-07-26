package main

import (
	"fmt"
	"time"
)

const (
	carbonMonoxide = "carbonMonoxide"
	airQuality     = "airQuality"
	raindrops      = "raindrops"
	soilMoisture   = "soilMoisture"
)

type Data struct {
	DataType  string
	Value     float32
	TimeStamp time.Time
}

type DataResponse struct {
	Data
	Category int
}

func (d Data) String() string {
	return fmt.Sprintf("DataType: %s, Value: %f, Time: %v", d.DataType, d.Value, d.TimeStamp)
}

func (dr DataResponse) String() string {
	return fmt.Sprintf("%s, Category: %d", dr.Data, dr.Category)
}
