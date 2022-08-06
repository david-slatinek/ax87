package server

import (
	"api/cache"
	"api/db"
	"api/model"
	pb "api/schema"
	"api/util"
	"context"
	"encoding/json"
	"errors"
	"log"
)

// Server is a struct that acts as an intermediate layer between db.DB and grpc default server.
type Server struct {
	pb.UnimplementedRequestServer
	// db.DB field.
	DBService *db.DB
	// Redis cache cache.
	cache *cache.Cache
	// Development mode.
	Development bool
}

// CreateCache creates server cache.
func (server *Server) CreateCache() {
	ca := cache.Cache{}
	ca.Load()
	err := ca.Create()

	if err == nil {
		server.cache = &ca
	} else {
		if server.Development {
			log.Printf("Invalid cache, error: %v", err)
		}
		server.cache = nil
	}
}

// Close db, redis connection.
func (server *Server) Close() {
	server.DBService.Close()

	if server.cache != nil {
		if err := server.cache.Close(); err != nil && server.Development {
			log.Printf("Error with cache.Close, error: %v", err)
		}
	}
}

// Add new data to the db.
func (server *Server) Add(_ context.Context, data *pb.Data) (*pb.Reply, error) {
	if data == nil {
		return &pb.Reply{}, errors.New("data can't be nil")
	}

	if data.GetDataType() == pb.DataType_NONE {
		return &pb.Reply{}, errors.New("invalid data type")
	}

	d := model.Data{
		DataType:  data.GetDataType().String(),
		Value:     data.GetValue(),
		Timestamp: data.GetTimestamp().AsTime(),
	}
	server.DBService.Add(&d)

	err := server.AddToCache(&model.DataResponse{
		Data:     d,
		Category: util.GetCategory(int(d.Value), d.DataType),
	})

	if err != nil && server.Development {
		log.Printf("Error with cache add, error: %v", err)
	}

	return &pb.Reply{}, nil
}

// AddToCache adds model.DataResponse object to the cache.
func (server *Server) AddToCache(dr *model.DataResponse) error {
	if server.cache != nil {
		return server.cache.Add(dr)
	}
	return errors.New("server.cache is nil")
}

// Latest returns the latest data for the requested dataType.
func (server *Server) Latest(_ context.Context, request *pb.DataRequest) (*pb.DataWithCategory, error) {
	if request == nil {
		return nil, errors.New("request can't be nil")
	}

	var latest *model.DataResponse
	var ok = false

	if server.cache != nil {
		value, err := server.cache.Get(request.DataType.String())

		if err != nil && server.Development {
			log.Printf("Error with cache get, error: %v", err)
		} else {
			if err := json.Unmarshal([]byte(value), &latest); err != nil && server.Development {
				log.Printf("Error with json.Unmarshal, error: %v", err)
			} else {
				ok = true
			}
		}
	}

	if !ok {
		var err error

		latest, err = server.DBService.Latest(request.GetDataType().String())
		if err != nil {
			return nil, err
		}

		if err = server.AddToCache(latest); err != nil && server.Development {
			log.Printf("Error with cache set, error: %v", err)
		}
	}

	return latest.Convert(), nil
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
