package main

type DataType int

const (
	carbonMonoxide DataType = iota
	airQuality
	raindrops
	soilMoisture
)

func (dt DataType) String() string {
	return []string{"carbon_monoxide", "air_quality", "raindrops", "soil_moisture"}[dt]
}

type Data struct {
	DataType DataType
	Value    float32
}
