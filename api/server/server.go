package server

import (
	"api/db"
	"api/model"
	pb "api/schema"
	"api/util"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v9"
	"log"
	"time"
)

// Server is a struct that acts as an intermediate layer between db.DB and grpc default server.
type Server struct {
	pb.UnimplementedRequestServer
	// db.DB field.
	DBService *db.DB
	// Redis cache client.
	client *redis.Client
}

// CreateClient creates server cache client.
func (server *Server) CreateClient() {
	server.client = redis.NewClient(&redis.Options{
		Addr:            "localhost:6379",
		Password:        "",
		DB:              0,
		MaxRetries:      5,
		MinRetryBackoff: time.Millisecond * 15,
	})

	_, err := server.client.Ping(context.Background()).Result()
	if err != nil {
		log.Printf("Invalid cache client, error: %v", err)
		server.client = nil
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

	if err != nil {
		log.Printf("Error with cache set, error: %v", err)
	}

	return &pb.Reply{}, nil
}

// AddToCache adds model.DataResponse object to cache.
func (server *Server) AddToCache(dr *model.DataResponse) error {
	if server.client != nil {
		return server.client.Set(context.Background(), dr.DataType, dr, time.Minute*5).Err()
	}
	return errors.New("server.client is nil")
}

// Latest returns the latest data for the requested dataType.
func (server *Server) Latest(_ context.Context, request *pb.DataRequest) (*pb.DataWithCategory, error) {
	if request == nil {
		return nil, errors.New("request can't be nil")
	}

	var latest *model.DataResponse
	var ok = false

	if server.client != nil {
		value, err := server.client.Get(context.Background(), request.DataType.String()).Result()

		if err != nil {
			log.Printf("Error with cache get, error: %v", err)
		} else {
			if err := json.Unmarshal([]byte(value), &latest); err != nil {
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

		if err = server.AddToCache(latest); err != nil {
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
