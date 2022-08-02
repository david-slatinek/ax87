package main

import (
	"api/env"
	model "api/models"
	pb "api/schema"
	"api/util"
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

// Test Server.Add.
func TestServer_Add(t *testing.T) {
	_ = env.Load("env/test.env")

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

// Test Server.Latest.
func TestServer_Latest(t *testing.T) {
	_ = env.Load("env/test.env")

	db := DB{}
	db.LoadFields()
	_ = db.Connect()
	defer db.client.Close()

	_ = db.Init()

	data := model.Data{
		DataType:  util.AirQuality,
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

	if !util.Equals(dc, expected) {
		t.Error("Objects are not the same")
		t.Errorf("Expected: %v", expected)
		t.Errorf("Result: %v", &dc)
		t.FailNow()
	}
}

// Test Server.Last24H.
func TestServer_Last24H(t *testing.T) {
	_ = env.Load("env/test.env")

	db := DB{}
	db.LoadFields()
	_ = db.Connect()
	defer db.client.Close()

	_ = db.Init()

	creationTime := time.Now().Round(0)

	var objects = []struct {
		data     model.Data
		expected pb.DataWithCategory
	}{
		{model.Data{
			DataType:  util.Raindrops,
			Value:     77,
			TimeStamp: creationTime,
		}, pb.DataWithCategory{
			Data: &pb.Data{
				DataType:  pb.DataType_RAINDROPS,
				Value:     77,
				Timestamp: timestamppb.New(creationTime),
			},
			Category: int32(util.GetCategory(77, util.Raindrops)),
		}},
		{model.Data{
			DataType:  util.Raindrops,
			Value:     87,
			TimeStamp: creationTime.Add(time.Second * -5),
		}, pb.DataWithCategory{
			Data: &pb.Data{
				DataType:  pb.DataType_RAINDROPS,
				Value:     87,
				Timestamp: timestamppb.New(creationTime.Add(time.Second * -5)),
			},
			Category: int32(util.GetCategory(87, util.Raindrops)),
		}},
		{model.Data{
			DataType:  util.Raindrops,
			Value:     224,
			TimeStamp: creationTime.Add(time.Minute * -10),
		}, pb.DataWithCategory{
			Data: &pb.Data{
				DataType:  pb.DataType_RAINDROPS,
				Value:     224,
				Timestamp: timestamppb.New(creationTime.Add(time.Minute * -10)),
			},
			Category: int32(util.GetCategory(224, util.Raindrops)),
		}},
		{model.Data{
			DataType:  util.Raindrops,
			Value:     400,
			TimeStamp: creationTime.Add(time.Hour * -7),
		}, pb.DataWithCategory{
			Data: &pb.Data{
				DataType:  pb.DataType_RAINDROPS,
				Value:     400,
				Timestamp: timestamppb.New(creationTime.Add(time.Hour * -7)),
			},
			Category: int32(util.GetCategory(400, util.Raindrops)),
		}},
		{model.Data{
			DataType:  util.Raindrops,
			Value:     21,
			TimeStamp: creationTime.Add(time.Hour * -12),
		}, pb.DataWithCategory{
			Data: &pb.Data{
				DataType:  pb.DataType_RAINDROPS,
				Value:     21,
				Timestamp: timestamppb.New(creationTime.Add(time.Hour * -12)),
			},
			Category: int32(util.GetCategory(21, util.Raindrops)),
		}},
	}

	for _, v := range objects {
		db.Add(&v.data)
	}

	server := Server{dbService: &db}

	dr, err := server.Last24H(context.Background(), &pb.DataRequest{DataType: pb.DataType_RAINDROPS})
	if err != nil {
		t.Fatalf("Expected nil with Last24H, got %v", err)
	}

	if len(dr.Data) != len(objects) {
		t.Fatalf("Expected length %d, got %d", len(objects), len(dr.Data))
	}

	for k, v := range objects {
		if !util.Equals(&v.expected, dr.Data[k]) {
			t.Error("Objects are not the same")
			t.Errorf("Expected: %v", &v.expected)
			t.Errorf("Result: %v", dr.Data[k])
		}
	}
}

func TestServer_Median(t *testing.T) {
	_ = env.Load("env/test.env")

	db := DB{}
	db.LoadFields()
	_ = db.Connect()
	defer db.client.Close()

	_ = db.Init()

	creationTime := time.Now().Round(0)

	var objects = [5]model.Data{
		{
			DataType:  util.SoilMoisture,
			Value:     222,
			TimeStamp: creationTime.Add(time.Minute * -3),
		},
		{
			DataType:  util.SoilMoisture,
			Value:     312,
			TimeStamp: creationTime.Add(time.Second * -4),
		},
		{
			DataType:  util.SoilMoisture,
			Value:     294,
			TimeStamp: creationTime.Add(time.Second * -5),
		},
		{
			DataType:  util.SoilMoisture,
			Value:     431,
			TimeStamp: creationTime.Add(time.Hour * -17),
		},
		{
			DataType:  util.SoilMoisture,
			Value:     401,
			TimeStamp: creationTime.Add(time.Hour * -1),
		},
	}

	for _, v := range objects {
		db.Add(&v)
	}

	server := Server{dbService: &db}

	dc, err := server.Median(context.Background(), &pb.DataRequest{DataType: pb.DataType_SOIL_MOISTURE})
	if err != nil {
		t.Fatalf("Expected nil with Median, got %v", err)
	}

	dr := pb.DataWithCategory{
		Data: &pb.Data{
			DataType:  pb.DataType_SOIL_MOISTURE,
			Value:     312,
			Timestamp: timestamppb.New(creationTime.Add(time.Second * -4)),
		},
		Category: int32(util.GetCategory(312, util.SoilMoisture)),
	}

	if !util.Equals(dc, &dr) {
		t.Error("Objects are not the same")
		t.Errorf("Expected: %v", &dr)
		t.Errorf("Result: %v", dc)
	}
}

// Test Server.Max.
func TestServer_Max(t *testing.T) {
	_ = env.Load("env/test.env")

	db := DB{}
	db.LoadFields()
	_ = db.Connect()
	defer db.client.Close()

	_ = db.Init()

	creationTime := time.Now().Round(0)

	var objects = [5]model.Data{
		{
			DataType:  util.CarbonMonoxide,
			Value:     65,
			TimeStamp: creationTime.Add(time.Minute * -2),
		},
		{
			DataType:  util.CarbonMonoxide,
			Value:     741,
			TimeStamp: creationTime.Add(time.Second * -45),
		},
		{
			DataType:  util.CarbonMonoxide,
			Value:     132,
			TimeStamp: creationTime.Add(time.Minute * -22),
		},
		{
			DataType:  util.CarbonMonoxide,
			Value:     387,
			TimeStamp: creationTime.Add(time.Hour * -3),
		},
		{
			DataType:  util.CarbonMonoxide,
			Value:     25,
			TimeStamp: creationTime.Add(time.Minute * -37),
		},
	}

	for _, v := range objects {
		db.Add(&v)
	}

	server := Server{dbService: &db}

	dc, err := server.Max(context.Background(), &pb.DataRequest{DataType: pb.DataType_CARBON_MONOXIDE})
	if err != nil {
		t.Fatalf("Expected nil with Max, got %v", err)
	}

	dr := pb.DataWithCategory{
		Data: &pb.Data{
			DataType:  pb.DataType_CARBON_MONOXIDE,
			Value:     741,
			Timestamp: timestamppb.New(creationTime.Add(time.Second * -45)),
		},
		Category: int32(util.GetCategory(741, util.CarbonMonoxide)),
	}

	if !util.Equals(dc, &dr) {
		t.Error("Objects are not the same")
		t.Errorf("Expected: %v", &dr)
		t.Errorf("Result: %v", dc)
	}
}

// Test Server.Min.
func TestServer_Min(t *testing.T) {
	_ = env.Load("env/test.env")

	db := DB{}
	db.LoadFields()
	_ = db.Connect()
	defer db.client.Close()

	_ = db.Init()

	creationTime := time.Now().Round(0)

	var objects = [5]model.Data{
		{
			DataType:  util.AirQuality,
			Value:     33,
			TimeStamp: creationTime.Add(time.Second * -2),
		},
		{
			DataType:  util.AirQuality,
			Value:     331,
			TimeStamp: creationTime.Add(time.Hour * -4),
		},
		{
			DataType:  util.AirQuality,
			Value:     363,
			TimeStamp: creationTime.Add(time.Minute * -35),
		},
		{
			DataType:  util.AirQuality,
			Value:     182,
			TimeStamp: creationTime.Add(time.Hour * -7),
		},
		{
			DataType:  util.AirQuality,
			Value:     167,
			TimeStamp: creationTime.Add(time.Second * -37),
		},
	}

	for _, v := range objects {
		db.Add(&v)
	}

	server := Server{dbService: &db}

	dc, err := server.Min(context.Background(), &pb.DataRequest{DataType: pb.DataType_AIR_QUALITY})
	if err != nil {
		t.Fatalf("Expected nil with Min, got %v", err)
	}

	dr := pb.DataWithCategory{
		Data: &pb.Data{
			DataType:  pb.DataType_AIR_QUALITY,
			Value:     33,
			Timestamp: timestamppb.New(creationTime.Add(time.Second * -2)),
		},
		Category: int32(util.GetCategory(33, util.AirQuality)),
	}

	if !util.Equals(dc, &dr) {
		t.Error("Objects are not the same")
		t.Errorf("Expected: %v", &dr)
		t.Errorf("Result: %v", dc)
	}
}
