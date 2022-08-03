package model

import (
	pb "api/schema"
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// DataResponse is a struct for serializing data from the database.
type DataResponse struct {
	// Embedded Data struct.
	Data
	// Data.Value category. Check functions util.GetCategory.
	Category int `json:"category"`
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
			Timestamp: timestamppb.New(dr.Timestamp),
		},
		Category: int32(dr.Category),
	}
}

// Equals compares two DataResponse structures.
func (dr *DataResponse) Equals(b *DataResponse) bool {
	return dr.DataType == b.DataType && dr.Value == b.Value && dr.Timestamp.Equal(b.Timestamp) && dr.Category == b.Category
}

func (dr DataResponse) MarshalBinary() ([]byte, error) {
	return json.Marshal(dr)
}
