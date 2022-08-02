package main

import (
	"api/env"
	pb "api/schema"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
)

func main() {
	err := env.Load("env/.env")
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = env.Validate()
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		_, cancel := context.WithTimeout(context.Background(), 15)
		defer cancel()
		grpcServer.GracefulStop()
		log.Println("Server is shutting down...")
		os.Exit(0)
	}()

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve, error: %v", err)
	}
}
