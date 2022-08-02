package model

import (
	pb "api/schema"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// DataResponse is a struct for serializing data from the database.
type DataResponse struct {
	// Embedded Data struct.
	Data
	// Data.Value category. Check functions MapCO2, MapAir, and MapValue.
	Category int
}

// String returns DataResponse fields in a string.
func (dr *DataResponse) String() string {
	return fmt.Sprintf("%v, Category: %d", dr.Data, dr.Category)
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

// Compare two DataResponse structures.
func (dr *DataResponse) Compare(b *DataResponse) bool {
	return dr.DataType == b.DataType && dr.Value == b.Value && dr.TimeStamp.Equal(b.TimeStamp) && dr.Category == b.Category
}
