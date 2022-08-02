package model

import (
	pb "api/schema"
	"api/util"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
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

// String returns Data fields in a string.
func (d *Data) String() string {
	return fmt.Sprintf("DataType: %s, Value: %f, Time: %v", d.DataType, d.Value, d.TimeStamp)
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
		Category: int32(util.GetCategory(int(d.Value), d.DataType)),
	}
}

// Compare two Data structures.
func (d *Data) Compare(b *Data) bool {
	return d.DataType == b.DataType && d.Value == b.Value && d.TimeStamp.Equal(b.TimeStamp)
}
