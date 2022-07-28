package main

import (
	pb "api/schema"
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedRequestServer
	// DB field.
	dbService *DB
}

func (server *Server) Add(_ context.Context, data *pb.Data) (*pb.EmptyReply, error) {
	d := Data{
		DataType:  data.GetDataType().String(),
		Value:     data.GetValue(),
		TimeStamp: data.GetTimestamp().AsTime(),
	}
	server.dbService.Add(&d)
	return &pb.EmptyReply{}, nil
}

func (server *Server) Latest(_ context.Context, request *pb.DataRequest) (*pb.DataWithCategory, error) {
	latest, err := server.dbService.Latest(request.GetDataType().String())

	if err != nil {
		return nil, err
	}

	dc := pb.DataWithCategory{
		Data: &pb.Data{
			DataType:  pb.DataType(pb.DataType_value[latest.DataType]),
			Value:     latest.Value,
			Timestamp: timestamppb.New(latest.TimeStamp),
		},
		Category: int32(latest.Category),
	}

	return &dc, err
}

func (server *Server) Last24H(_ context.Context, request *pb.DataRequest) (*pb.DataRepeated, error) {
	last, err := server.dbService.Last24H(request.GetDataType().String())

	if err != nil {
		return nil, err
	}

	var dc []*pb.DataWithCategory

	for _, value := range *last {
		dc = append(dc, &pb.DataWithCategory{
			Data: &pb.Data{
				DataType:  pb.DataType(pb.DataType_value[value.DataType]),
				Value:     value.Value,
				Timestamp: timestamppb.New(value.TimeStamp),
			},
			Category: int32(value.Category),
		})

	}

	return &pb.DataRepeated{Data: dc}, nil
}
