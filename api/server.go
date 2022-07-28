package main

import (
	pb "api/schema"
	"context"
	"errors"
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

func (server *Server) Latest(_ context.Context, request *pb.DataRequest) (*pb.DataWithCategory, error) {
	if request == nil {
		return nil, errors.New("request can't be nil")
	}

	latest, err := server.dbService.Latest(request.GetDataType().String())

	if err != nil {
		return nil, err
	}

	return latest.Convert(), err
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
		dc = append(dc, value.Convert())
	}

	return &pb.DataRepeated{Data: dc}, nil
}

func (server *Server) Median(_ context.Context, request *pb.DataRequest) (*pb.Data, error) {
	if request == nil {
		return nil, errors.New("request can't be nil")
	}

	median, err := server.dbService.Median(request.GetDataType().String())
	if err != nil {
		return nil, err
	}

	return median.Convert(), nil
}

func (server *Server) Max(_ context.Context, request *pb.DataRequest) (*pb.Data, error) {
	if request == nil {
		return nil, errors.New("request can't be nil")
	}

	max, err := server.dbService.Max(request.GetDataType().String())
	if err != nil {
		return nil, err
	}

	return max.Convert(), nil
}

func (server *Server) Min(_ context.Context, request *pb.DataRequest) (*pb.Data, error) {
	if request == nil {
		return nil, errors.New("request can't be nil")
	}

	min, err := server.dbService.Min(request.GetDataType().String())
	if err != nil {
		return nil, err
	}

	return min.Convert(), nil
}
