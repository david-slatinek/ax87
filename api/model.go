package main

type DataType int

const (
	carbonMonoxide DataType = iota
	airQuality
	raindrops
	soilMoisture
)

func (dt DataType) String() string {
	return []string{"carbonMonoxide", "airQuality", "raindrops", "soilMoisture"}[dt]
}

type Data struct {
	DataType DataType
	Value    float32
}
