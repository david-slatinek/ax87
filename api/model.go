package main

type DataType int

const (
	carbon_monoxide DataType = iota
	air_quality
	raindrops
	soil_moisture
)

type Data struct {
	DataType DataType
	Value    float32
}
