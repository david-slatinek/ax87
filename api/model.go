package main

import (
	pb "api/schema"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

const (
	// Carbon monoxide constant.
	carbonMonoxide = "CARBON_MONOXIDE"
	// Air quality constant.
	airQuality = "AIR_QUALITY"
	// Raindrops constant.
	raindrops = "RAINDROPS"
	// Soil moisture constant.
	soilMoisture = "SOIL_MOISTURE"
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
	return fmt.Sprintf("%v, Category: %d", dr.Data, dr.Category)
}

// Convert Data to pb.Data.
func (d *Data) Convert() *pb.Data {
	return &pb.Data{
		DataType:  pb.DataType(pb.DataType_value[d.DataType]),
		Value:     d.Value,
		Timestamp: timestamppb.New(d.TimeStamp),
	}
}

// ConvertToDC converts Data to pb.DataWithCategory.
func (d *Data) ConvertToDC() *pb.DataWithCategory {
	return &pb.DataWithCategory{
		Data:     d.Convert(),
		Category: int32(GetCategory(int(d.Value), d.DataType)),
	}
}

// Convert DataResponse to pb.DataWithCategory.
func (dr *DataResponse) Convert() *pb.DataWithCategory {
	return &pb.DataWithCategory{
		Data: &pb.Data{
			DataType:  pb.DataType(pb.DataType_value[dr.DataType]),
			Value:     dr.Value,
			Timestamp: timestamppb.New(dr.TimeStamp),
		},
		Category: int32(dr.Category),
	}
}

// Compare two DataResponse structures. Compares all fields.
func (dr *DataResponse) Compare(b *DataResponse) bool {
	return dr.DataType == b.DataType && dr.Value == b.Value && dr.TimeStamp.Equal(b.TimeStamp) && dr.Category == b.Category
}

func (d *Data) Compare(b *Data) bool {
	return d.DataType == b.DataType && d.Value == b.Value && d.TimeStamp.Equal(b.TimeStamp)
}

// Equals compares two pb.DataWithCategory structures. Compares all fields.
func Equals(a, b *pb.DataWithCategory) bool {
	return a.Data.DataType == b.Data.DataType && a.Data.Value == b.Data.Value &&
		a.Data.Timestamp.AsTime().Equal(b.Data.Timestamp.AsTime()) && a.Category == b.Category
}
