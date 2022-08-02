package server

import (
	"api/db"
	"api/model"
	pb "api/schema"
	"context"
	"errors"
)

// Server is a struct that acts as an intermediate layer between db.DB and grpc default server.
type Server struct {
	pb.UnimplementedRequestServer
	// db.DB field.
	DBService *db.DB
}

// Add new data to the db.
func (server *Server) Add(_ context.Context, data *pb.Data) (*pb.Reply, error) {
	if data == nil {
		return &pb.Reply{}, errors.New("data can't be nil")
	}

	d := model.Data{
		DataType:  data.GetDataType().String(),
		Value:     data.GetValue(),
		Timestamp: data.GetTimestamp().AsTime(),
	}
	server.DBService.Add(&d)
	return &pb.Reply{}, nil
}

// Latest returns the latest data for the requested dataType.
func (server *Server) Latest(_ context.Context, request *pb.DataRequest) (*pb.DataWithCategory, error) {
	if request == nil {
		return nil, errors.New("request can't be nil")
	}

	latest, err := server.DBService.Latest(request.GetDataType().String())

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

	last, err := server.DBService.Last24H(request.GetDataType().String())
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

	median, err := server.DBService.Median(request.GetDataType().String())
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

	max, err := server.DBService.Max(request.GetDataType().String())
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

	min, err := server.DBService.Min(request.GetDataType().String())
	if err != nil {
		return nil, err
	}

	return min.Convert(), nil
}
