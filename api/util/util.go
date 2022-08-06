package util

import (
	pb "api/schema"
	"context"
	"crypto/tls"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"math"
	"os"
	"time"
)

var rateLimiter = rate.NewLimiter(rate.Every(time.Second), 10)

const (
	// CarbonMonoxide constant.
	CarbonMonoxide = "CARBON_MONOXIDE"
	// AirQuality constant.
	AirQuality = "AIR_QUALITY"
	// Raindrops constant.
	Raindrops = "RAINDROPS"
	// SoilMoisture constant.
	SoilMoisture = "SOIL_MOISTURE"

	// EnvTestFilePath location for test.env.
	EnvTestFilePath = "../env/test.env"
)

// MapCO2 value to 7 categories, with 1 being the best.
func MapCO2(value int) int {
	if value <= 30 {
		return 1
	} else if value > 30 && value <= 70 {
		return 2
	} else if value > 70 && value <= 150 {
		return 3
	} else if value > 150 && value <= 200 {
		return 4
	} else if value > 200 && value <= 400 {
		return 5
	} else if value > 400 && value <= 800 {
		return 6
	}
	return 7
}

// MapAir quality value to 6 categories, with 1 being the best.
func MapAir(value int) int {
	if value <= 50 {
		return 1
	} else if value > 50 && value <= 100 {
		return 2
	} else if value > 100 && value <= 150 {
		return 3
	} else if value > 150 && value <= 200 {
		return 4
	} else if value > 200 && value <= 300 {
		return 5
	}
	return 6
}

// MapValue from one range to another.
func MapValue(x float64, inMin float64, inMax float64, outMin float64, outMax float64) int {
	return int((math.Round(x-inMin)*(outMax-outMin)/(inMax-inMin) + outMin) + 0.5)
}

// equalsData compares two pb.DataWithCategory structures. Checks only pb.Data field.
func equalsData(a, b *pb.DataWithCategory) bool {
	return a.Data.DataType == b.Data.DataType && a.Data.Value == b.Data.Value &&
		a.Data.Timestamp.AsTime().Equal(b.Data.Timestamp.AsTime())
}

// Equals compares two pb.DataWithCategory structures.
func Equals(a, b *pb.DataWithCategory) bool {
	return equalsData(a, b) && a.Category == b.Category
}

// GetCategory returns category for a specific dataType.
func GetCategory(value int, dataType string) int {
	switch dataType {
	case CarbonMonoxide:
		return MapCO2(value)
	case AirQuality:
		return MapAir(value)
	case Raindrops:
		return MapValue(float64(value), 0, 1024, 1, 4)
	case SoilMoisture:
		return MapValue(float64(value), 489, 238, 0, 100)
	}
	return -1
}

func RateLimit(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if err := rateLimiter.Wait(ctx); err != nil && os.Getenv("GO_ENV") == "development" {
		log.Println("Interceptor error:", err)
	}
	return handler(ctx, req)
}

func LoadTLS() (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair("cert/server-cert.pem", "cert/server-key.pem")
	if err != nil {
		return nil, err
	}

	config := tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	return credentials.NewTLS(&config), nil
}
