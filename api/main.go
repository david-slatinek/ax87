package main

import (
	pb "api/schema"
	"google.golang.org/grpc"
	"log"
	"net"
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

	db := DB{}
	db.LoadFields()

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
	server := Server{dbService: &db}
	pb.RegisterRequestServer(grpcServer, &server)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve, error: %v", err)
	}
}
