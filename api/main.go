package main

import (
	pb "api/schema"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	err := Load(".env")
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = Validate()
	if err != nil {
		log.Fatalf(err.Error())
	}

	db := DB{url: os.Getenv("INFLUXDB_URL"), token: os.Getenv("INFLUXDB_TOKEN"),
		org: os.Getenv("ORGANIZATION"), bucket: os.Getenv("BUCKET")}
	err = db.Connect()
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer db.client.Close()

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000 with error: %v", err)
	}

	grpcServer := grpc.NewServer()
	server := Server{DbService: &db}
	pb.RegisterRequestServer(grpcServer, &server)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve, error: %v", err)
	}
}
