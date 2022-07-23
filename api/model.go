package main

import "fmt"

const (
	carbonMonoxide = "carbonMonoxide"
	airQuality     = "airQuality"
	raindrops      = "raindrops"
	soilMoisture   = "soilMoisture"
)

type Data struct {
	DataType string
	Value    float32
}

type DataResponse struct {
	Data
	Category int
}

func (d Data) String() string {
	return fmt.Sprintf("DataType: %s, Value: %f", d.DataType, d.Value)
}

func (dr DataResponse) String() string {
	return fmt.Sprintf("%s, Category: %d", dr.Data, dr.Category)
}
