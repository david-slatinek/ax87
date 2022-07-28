package main

import (
	pb "api/schema"
	"context"
	"errors"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedRequestServer
	// DB field.
	dbService *DB
}

func (server *Server) Add(_ context.Context, data *pb.Data) (*pb.EmptyReply, error) {
	if data == nil {
		return &pb.EmptyReply{}, errors.New("data can't be nil")
	}

	d := Data{
		DataType:  data.GetDataType().String(),
		Value:     data.GetValue(),
		TimeStamp: data.GetTimestamp().AsTime(),
	}
	server.dbService.Add(&d)
	return &pb.EmptyReply{}, nil
}

func Convert(t *DataResponse) *pb.DataWithCategory {
	return &pb.DataWithCategory{
		Data: &pb.Data{
			DataType:  pb.DataType(pb.DataType_value[t.DataType]),
			Value:     t.Value,
			Timestamp: timestamppb.New(t.TimeStamp),
		},
		Category: int32(t.Category),
	}
}

func (server *Server) Latest(_ context.Context, request *pb.DataRequest) (*pb.DataWithCategory, error) {
	if request == nil {
		return nil, errors.New("request can't be nil")
	}

	latest, err := server.dbService.Latest(request.GetDataType().String())

	if err != nil {
		return nil, err
	}

	return Convert(latest), err
}

func (server *Server) Last24H(_ context.Context, request *pb.DataRequest) (*pb.DataRepeated, error) {
	if request == nil {
		return nil, errors.New("request can't be nil")
	}

	last, err := server.dbService.Last24H(request.GetDataType().String())
	if err != nil {
		return nil, err
	}

	var dc []*pb.DataWithCategory

	for _, value := range *last {
		dc = append(dc, Convert(&value))
	}

	return &pb.DataRepeated{Data: dc}, nil
}

func ConvertData(data *Data) *pb.Data {
	return &pb.Data{
		DataType:  pb.DataType(pb.DataType_value[data.DataType]),
		Value:     data.Value,
		Timestamp: timestamppb.New(data.TimeStamp),
	}
}

func (server *Server) Median(_ context.Context, request *pb.DataRequest) (*pb.Data, error) {
	if request == nil {
		return nil, errors.New("request can't be nil")
	}

	median, err := server.dbService.Median(request.GetDataType().String())
	if err != nil {
		return nil, err
	}

	return ConvertData(median), nil
}
