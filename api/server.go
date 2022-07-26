package main

import (
	pb "api/schema"
	"context"
)

type Server struct {
	pb.UnimplementedRequestServer
	Db *DB
}

func (s *Server) Add(_ context.Context, data *pb.Data) (*pb.Reply, error) {
	d := Data{
		DataType:  data.GetDataType().String(),
		Value:     data.GetValue(),
		TimeStamp: data.GetTimestamp().AsTime(),
	}
	s.Db.Add(&d)
	return &pb.Reply{}, nil
}
