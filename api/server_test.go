package main

import (
	pb "api/schema"
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestServer_Add(t *testing.T) {
	_ = Load("test.env")

	db := DB{}
	db.LoadFields()
	_ = db.Connect()
	defer db.client.Close()

	_ = db.Init()

	server := Server{dbService: &db}

	d := pb.Data{
		DataType:  pb.DataType_CARBON_MONOXIDE,
		Value:     55,
		Timestamp: timestamppb.New(time.Now()),
	}

	reply, err := server.Add(context.Background(), &d)

	if err != nil {
		t.Fatalf("Expected nil with Server.Add, got %v", err)
	}

	r := pb.Reply{}

	if reply.String() != r.String() {
		t.Error("Objects not the same")
		t.Errorf("Expected: %v", &r)
		t.Errorf("Result: %v", reply)
		t.FailNow()
	}
}

func TestServer_Latest(t *testing.T) {
	_ = Load("test.env")

	db := DB{}
	db.LoadFields()
	_ = db.Connect()
	defer db.client.Close()

	_ = db.Init()

	data := Data{
		DataType:  airQuality,
		Value:     183,
		TimeStamp: time.Now().Round(0),
	}
	db.Add(&data)

	server := Server{dbService: &db}
	dc, err := server.Latest(context.Background(), &pb.DataRequest{DataType: pb.DataType_AIR_QUALITY})

	if err != nil {
		t.Fatalf("Expected nil with Server.Latest, got %v", err)
	}

	expected := data.ConvertToDC()

	if !Equals(dc, expected) {
		t.Error("Objects are not the same")
		t.Errorf("Expected: %v", expected)
		t.Errorf("Result: %v", &dc)
		t.FailNow()
	}
}
