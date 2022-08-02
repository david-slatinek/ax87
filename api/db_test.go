package main

import (
	"api/env"
	"api/models"
	"api/util"
	"testing"
	"time"
)

// Test DB.LoadFields.
func TestDB_LoadFields(t *testing.T) {
	_ = env.Load("env/test.env")

	db := DB{}
	db.LoadFields()

	if db.token == "" {
		t.Error("Empty db.token")
	}

	if db.url == "" {
		t.Error("Empty db.url")
	}

	if db.org == "" {
		t.Error("Empty db.org")
	}

	if db.bucket == "" {
		t.Error("Empty db.bucket")
	}
}

// Test DB.Connect.
func TestDB_Connect(t *testing.T) {
	_ = env.Load("env/test.env")

	db := DB{}
	db.LoadFields()

	if err := db.Connect(); err != nil {
		t.Fatalf("Expected nil with Connect, got %v", err)
	}
	defer db.client.Close()
}

// Test DB.Init
func TestDB_Init(t *testing.T) {
	_ = env.Load("env/test.env")

	db := DB{}
	db.LoadFields()
	_ = db.Connect()
	defer db.client.Close()

	if err := db.Init(); err != nil {
		t.Fatalf("Expected nil with Init, got %v", err)
	}
}

// Test DB.Latest.
func TestDB_Latest(t *testing.T) {
	_ = env.Load("env/test.env")

	db := DB{}
	db.LoadFields()
	_ = db.Connect()
	defer db.client.Close()

	_ = db.Init()

	creationTime := time.Now().Round(0)

	var objects = []struct {
		data     model.Data
		expected model.DataResponse
	}{
		{model.Data{
			DataType:  util.CarbonMonoxide,
			Value:     45,
			TimeStamp: creationTime,
		}, model.DataResponse{
			Data: model.Data{
				DataType:  util.CarbonMonoxide,
				Value:     45,
				TimeStamp: creationTime,
			},
			Category: util.MapCO2(45),
		}},
		{model.Data{
			DataType:  util.AirQuality,
			Value:     125,
			TimeStamp: creationTime,
		}, model.DataResponse{
			Data: model.Data{
				DataType:  util.AirQuality,
				Value:     125,
				TimeStamp: creationTime,
			},
			Category: util.MapAir(125),
		}},
		{model.Data{
			DataType:  util.Raindrops,
			Value:     800,
			TimeStamp: creationTime,
		}, model.DataResponse{
			Data: model.Data{
				DataType:  util.Raindrops,
				Value:     800,
				TimeStamp: creationTime,
			},
			Category: util.MapValue(800, 0, 1024, 1, 4),
		}},
		{model.Data{
			DataType:  util.SoilMoisture,
			Value:     400,
			TimeStamp: creationTime,
		}, model.DataResponse{
			Data: model.Data{
				DataType:  util.SoilMoisture,
				Value:     400,
				TimeStamp: creationTime,
			},
			Category: util.MapValue(400, 489, 238, 0, 100),
		}},
	}

	for _, v := range objects {
		db.Add(&v.data)
	}

	types := [4]string{util.CarbonMonoxide, util.AirQuality, util.Raindrops, util.SoilMoisture}

	for k, v := range types {
		dr, err := db.Latest(v)
		if err != nil {
			t.Fatalf("Expected nil with Latest, got %v", err)
		}

		if !dr.Compare(&objects[k].expected) {
			t.Error("Objects are not the same")
			t.Errorf("Expected: %v", objects[k].expected)
			t.Errorf("Result: %v", dr)
		}
	}
}

// Test DB.Last24H.
func TestDB_Last24H(t *testing.T) {
	_ = env.Load("env/test.env")

	db := DB{}
	db.LoadFields()
	_ = db.Connect()
	defer db.client.Close()

	_ = db.Init()

	creationTime := time.Now().Round(0)

	var objects = []struct {
		data     model.Data
		expected model.DataResponse
	}{
		{model.Data{
			DataType:  util.CarbonMonoxide,
			Value:     250,
			TimeStamp: creationTime,
		}, model.DataResponse{
			Data: model.Data{
				DataType:  util.CarbonMonoxide,
				Value:     250,
				TimeStamp: creationTime,
			},
			Category: util.MapCO2(250),
		}},
		{model.Data{
			DataType:  util.CarbonMonoxide,
			Value:     55,
			TimeStamp: creationTime.Add(time.Second * -1),
		}, model.DataResponse{
			Data: model.Data{
				DataType:  util.CarbonMonoxide,
				Value:     55,
				TimeStamp: creationTime.Add(time.Second * -1),
			},
			Category: util.MapCO2(55),
		}},
		{model.Data{
			DataType:  util.CarbonMonoxide,
			Value:     420,
			TimeStamp: creationTime.Add(time.Second * -10),
		}, model.DataResponse{
			Data: model.Data{
				DataType:  util.CarbonMonoxide,
				Value:     420,
				TimeStamp: creationTime.Add(time.Second * -10),
			},
			Category: util.MapCO2(420),
		}},
		{model.Data{
			DataType:  util.CarbonMonoxide,
			Value:     69,
			TimeStamp: creationTime.Add(time.Minute * -1),
		}, model.DataResponse{
			Data: model.Data{
				DataType:  util.CarbonMonoxide,
				Value:     69,
				TimeStamp: creationTime.Add(time.Minute * -1),
			},
			Category: util.MapCO2(69),
		}},
		{model.Data{
			DataType:  util.CarbonMonoxide,
			Value:     170,
			TimeStamp: creationTime.Add(time.Minute * -2),
		}, model.DataResponse{
			Data: model.Data{
				DataType:  util.CarbonMonoxide,
				Value:     170,
				TimeStamp: creationTime.Add(time.Minute * -2),
			},
			Category: util.MapCO2(170),
		}},
	}

	for _, v := range objects {
		db.Add(&v.data)
	}

	dr, err := db.Last24H(util.CarbonMonoxide)
	if err != nil {
		t.Fatalf("Expected nil with Last24H, got %v", err)
	}

	if len(*dr) != len(objects) {
		t.Fatalf("Expected length %d, got %d", len(objects), len(*dr))
	}

	for k, v := range *dr {
		if !v.Compare(&objects[k].expected) {
			t.Error("Objects are not the same")
			t.Errorf("Expected: %v", objects[k].expected)
			t.Errorf("Result: %v", v)
		}
	}
}

// Test DB.Median.
func TestDB_Median(t *testing.T) {
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
			Value:     140,
			TimeStamp: creationTime.Add(time.Minute * -2),
		},
		{
			DataType:  util.AirQuality,
			Value:     205,
			TimeStamp: creationTime.Add(time.Second * -2),
		},
		{
			DataType:  util.AirQuality,
			Value:     270,
			TimeStamp: creationTime.Add(time.Second * -5),
		},
		{
			DataType:  util.AirQuality,
			Value:     33,
			TimeStamp: creationTime.Add(time.Second * -21),
		},
		{
			DataType:  util.AirQuality,
			Value:     195,
			TimeStamp: creationTime.Add(time.Minute * -1),
		},
	}

	for _, v := range objects {
		db.Add(&v)
	}

	d, err := db.Median(util.AirQuality)
	if err != nil {
		t.Fatalf("Expected nil with Median, got %v", err)
	}

	dr := model.DataResponse{
		Data: model.Data{
			DataType:  util.AirQuality,
			Value:     195,
			TimeStamp: creationTime.Add(time.Minute * -1),
		},
		Category: util.MapAir(195),
	}

	if !d.Compare(&dr) {
		t.Error("Objects are not the same")
		t.Errorf("Expected: %v", dr)
		t.Errorf("Result: %v", d)
	}
}

// Test DB.Max.
func TestDB_Max(t *testing.T) {
	_ = env.Load("env/test.env")

	db := DB{}
	db.LoadFields()
	_ = db.Connect()
	defer db.client.Close()

	_ = db.Init()

	creationTime := time.Now().Round(0)

	var objects = [5]model.Data{
		{
			DataType:  util.Raindrops,
			Value:     230,
			TimeStamp: creationTime.Add(time.Minute * -1),
		},
		{
			DataType:  util.Raindrops,
			Value:     420,
			TimeStamp: creationTime.Add(time.Second * -3),
		},
		{
			DataType:  util.Raindrops,
			Value:     114,
			TimeStamp: creationTime.Add(time.Second * -7),
		},
		{
			DataType:  util.Raindrops,
			Value:     47,
			TimeStamp: creationTime.Add(time.Second * -41),
		},
		{
			DataType:  util.Raindrops,
			Value:     842,
			TimeStamp: creationTime.Add(time.Minute * -2),
		},
	}

	for _, v := range objects {
		db.Add(&v)
	}

	d, err := db.Max(util.Raindrops)
	if err != nil {
		t.Fatalf("Expected nil with Max, got %v", err)
	}

	dr := model.DataResponse{
		Data: model.Data{
			DataType:  util.Raindrops,
			Value:     842,
			TimeStamp: creationTime.Add(time.Minute * -2),
		},
		Category: util.MapValue(842, 0, 1024, 1, 4),
	}

	if !d.Compare(&dr) {
		t.Error("Objects are not the same")
		t.Errorf("Expected: %v", dr)
		t.Errorf("Result: %v", d)
	}
}

// Test DB.Min.
func TestDB_Min(t *testing.T) {
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
			Value:     250,
			TimeStamp: creationTime.Add(time.Second * -10),
		},
		{
			DataType:  util.SoilMoisture,
			Value:     307,
			TimeStamp: creationTime.Add(time.Second * -17),
		},
		{
			DataType:  util.SoilMoisture,
			Value:     348,
			TimeStamp: creationTime.Add(time.Second * -2),
		},
		{
			DataType:  util.SoilMoisture,
			Value:     412,
			TimeStamp: creationTime.Add(time.Minute * -5),
		},
		{
			DataType:  util.SoilMoisture,
			Value:     440,
			TimeStamp: creationTime.Add(time.Minute * -1),
		},
	}

	for _, v := range objects {
		db.Add(&v)
	}

	d, err := db.Min(util.SoilMoisture)
	if err != nil {
		t.Fatalf("Expected nil with Min, got %v", err)
	}

	dr := model.DataResponse{
		Data: model.Data{
			DataType:  util.SoilMoisture,
			Value:     250,
			TimeStamp: creationTime.Add(time.Second * -10),
		},
		Category: util.MapValue(250, 489, 238, 0, 100),
	}

	if !d.Compare(&dr) {
		t.Error("Objects are not the same")
		t.Errorf("Expected: %v", dr)
		t.Errorf("Result: %v", d)
	}
}
