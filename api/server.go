package main

import (
	pb "api/schema"
	"context"
)

type Server struct {
	pb.UnimplementedRequestServer
	// DB field.
	DbService *DB
}

func (server *Server) Add(_ context.Context, data *pb.Data) (*pb.EmptyReply, error) {
	d := Data{
		DataType:  data.GetDataType().String(),
		Value:     data.GetValue(),
		TimeStamp: data.GetTimestamp().AsTime(),
	}
	server.DbService.Add(&d)
	return &pb.EmptyReply{}, nil
}
