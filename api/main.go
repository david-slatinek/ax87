package main

import (
	"api/db"
	"api/env"
	pb "api/schema"
	"api/server"
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

	dbb := db.DB{}
	dbb.LoadFields()

	err = dbb.Connect()
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer dbb.Close()

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000 with error: %v", err)
	}

	grpcServer := grpc.NewServer()
	srv := server.Server{DBService: &dbb}
	srv.CreateClient()
	pb.RegisterRequestServer(grpcServer, &srv)

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
