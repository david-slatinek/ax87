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

// Add new data to the db.
func (server *Server) Add(_ context.Context, data *pb.Data) (*pb.Reply, error) {
	if data == nil {
		return &pb.Reply{}, errors.New("data can't be nil")
	}

	d := Data{
		DataType:  data.GetDataType().String(),
		Value:     data.GetValue(),
		TimeStamp: data.GetTimestamp().AsTime(),
	}
	server.dbService.Add(&d)
	return &pb.Reply{}, nil
}

// Latest returns the latest data for the requested dataType.
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

// Last24H returns data for the last 24 hours for the requested dataType.
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

// Median returns median data for the last 24 hours for the requested dataType.
func (server *Server) Median(_ context.Context, request *pb.DataRequest) (*pb.DataWithCategory, error) {
	if request == nil {
		return nil, errors.New("request can't be nil")
	}

	median, err := server.dbService.Median(request.GetDataType().String())
	if err != nil {
		return nil, err
	}

	return median.Convert(), nil
}

// Max returns maximum data for the last 24 hours for the requested dataType.
func (server *Server) Max(_ context.Context, request *pb.DataRequest) (*pb.DataWithCategory, error) {
	if request == nil {
		return nil, errors.New("request can't be nil")
	}

	max, err := server.dbService.Max(request.GetDataType().String())
	if err != nil {
		return nil, err
	}

	return max.Convert(), nil
}

// Min returns minimum data for the last 24 hours for the requested dataType.
func (server *Server) Min(_ context.Context, request *pb.DataRequest) (*pb.DataWithCategory, error) {
	if request == nil {
		return nil, errors.New("request can't be nil")
	}

	min, err := server.dbService.Min(request.GetDataType().String())
	if err != nil {
		return nil, err
	}

	return min.Convert(), nil
}
