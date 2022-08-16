package main

import (
	"api/db"
	"api/env"
	pb "api/schema"
	"api/server"
	"api/util"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	dbb := &db.DB{}
	dbb.Load()

	err = dbb.Connect()
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer dbb.Close()

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000 with error: %v", err)
	}
	defer func(listener net.Listener) {
		if err := listener.Close(); err != nil {
			log.Printf("Expected nil with listener.Close, got %v", err)
		}
	}(listener)

	tlsCred, err := util.LoadTLS()
	if err != nil {
		log.Fatalf("Failed to use TLS, error: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(util.RateLimit),
		grpc.Creds(tlsCred),
	)
	srv := &server.Server{
		DBService:   dbb,
		Development: os.Getenv("GO_ENV") == "development",
	}
	srv.CreateCache()
	pb.RegisterRequestServer(grpcServer, srv)

	if srv.Development {
		reflection.Register(grpcServer)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		_, cancel := context.WithTimeout(context.Background(), 15)
		defer cancel()

		grpcServer.GracefulStop()
		srv.Close()

		log.Println("Server is shutting down...")
		os.Exit(0)
	}()

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve, error: %v", err)
	}
}
