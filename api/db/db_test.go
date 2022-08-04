package db_test

import (
	"api/db"
	"api/env"
	"api/model"
	"api/util"
	"testing"
	"time"
)

// Test DB.Connect.
func TestDB_Connect(t *testing.T) {
	_ = env.Load(util.EnvTestFilePath)

	dbb := db.DB{}
	dbb.Load()

	if err := dbb.Connect(); err != nil {
		t.Fatalf("Expected nil with Connect, got %v", err)
	}
	defer dbb.Close()
}

// Test DB.Init.
func TestDB_Init(t *testing.T) {
	_ = env.Load(util.EnvTestFilePath)

	dbb := db.DB{}
	dbb.Load()
	_ = dbb.Connect()
	defer dbb.Close()

	if err := dbb.Init(); err != nil {
		t.Fatalf("Expected nil with Init, got %v", err)
	}
}

// Test DB.Latest.
func TestDB_Latest(t *testing.T) {
	_ = env.Load(util.EnvTestFilePath)

	dbb := db.DB{}
	dbb.Load()
	_ = dbb.Connect()
	defer dbb.Close()

	_ = dbb.Init()

	creationTime := time.Now().Round(0)

	var objects = []struct {
		data     model.Data
		expected model.DataResponse
	}{
		{model.Data{
			DataType:  util.CarbonMonoxide,
			Value:     45,
			Timestamp: creationTime,
		}, model.DataResponse{
			Data: model.Data{
				DataType:  util.CarbonMonoxide,
				Value:     45,
				Timestamp: creationTime,
			},
			Category: util.GetCategory(45, util.CarbonMonoxide),
		}},
		{model.Data{
			DataType:  util.AirQuality,
			Value:     125,
			Timestamp: creationTime,
		}, model.DataResponse{
			Data: model.Data{
				DataType:  util.AirQuality,
				Value:     125,
				Timestamp: creationTime,
			},
			Category: util.GetCategory(125, util.AirQuality),
		}},
		{model.Data{
			DataType:  util.Raindrops,
			Value:     800,
			Timestamp: creationTime,
		}, model.DataResponse{
			Data: model.Data{
				DataType:  util.Raindrops,
				Value:     800,
				Timestamp: creationTime,
			},
			Category: util.GetCategory(800, util.Raindrops),
		}},
		{model.Data{
			DataType:  util.SoilMoisture,
			Value:     400,
			Timestamp: creationTime,
		}, model.DataResponse{
			Data: model.Data{
				DataType:  util.SoilMoisture,
				Value:     400,
				Timestamp: creationTime,
			},
			Category: util.GetCategory(400, util.SoilMoisture),
		}},
	}

	for _, v := range objects {
		dbb.Add(&v.data)
	}

	types := [4]string{util.CarbonMonoxide, util.AirQuality, util.Raindrops, util.SoilMoisture}

	for k, v := range types {
		dr, err := dbb.Latest(v)
		if err != nil {
			t.Fatalf("Expected nil with Latest, got %v", err)
		}

		if !dr.Equals(&objects[k].expected) {
			t.Error("Objects are not the same")
			t.Errorf("Expected: %v", objects[k].expected)
			t.Errorf("Result: %v", dr)
		}
	}
}

// Test DB.Last24H.
func TestDB_Last24H(t *testing.T) {
	_ = env.Load(util.EnvTestFilePath)

	dbb := db.DB{}
	dbb.Load()
	_ = dbb.Connect()
	defer dbb.Close()

	_ = dbb.Init()

	creationTime := time.Now().Round(0)

	var objects = []struct {
		data     model.Data
		expected model.DataResponse
	}{
		{model.Data{
			DataType:  util.CarbonMonoxide,
			Value:     250,
			Timestamp: creationTime,
		}, model.DataResponse{
			Data: model.Data{
				DataType:  util.CarbonMonoxide,
				Value:     250,
				Timestamp: creationTime,
			},
			Category: util.GetCategory(250, util.CarbonMonoxide),
		}},
		{model.Data{
			DataType:  util.CarbonMonoxide,
			Value:     55,
			Timestamp: creationTime.Add(time.Second * -1),
		}, model.DataResponse{
			Data: model.Data{
				DataType:  util.CarbonMonoxide,
				Value:     55,
				Timestamp: creationTime.Add(time.Second * -1),
			},
			Category: util.GetCategory(55, util.CarbonMonoxide),
		}},
		{model.Data{
			DataType:  util.CarbonMonoxide,
			Value:     420,
			Timestamp: creationTime.Add(time.Second * -10),
		}, model.DataResponse{
			Data: model.Data{
				DataType:  util.CarbonMonoxide,
				Value:     420,
				Timestamp: creationTime.Add(time.Second * -10),
			},
			Category: util.GetCategory(420, util.CarbonMonoxide),
		}},
		{model.Data{
			DataType:  util.CarbonMonoxide,
			Value:     69,
			Timestamp: creationTime.Add(time.Minute * -1),
		}, model.DataResponse{
			Data: model.Data{
				DataType:  util.CarbonMonoxide,
				Value:     69,
				Timestamp: creationTime.Add(time.Minute * -1),
			},
			Category: util.GetCategory(69, util.CarbonMonoxide),
		}},
		{model.Data{
			DataType:  util.CarbonMonoxide,
			Value:     170,
			Timestamp: creationTime.Add(time.Minute * -2),
		}, model.DataResponse{
			Data: model.Data{
				DataType:  util.CarbonMonoxide,
				Value:     170,
				Timestamp: creationTime.Add(time.Minute * -2),
			},
			Category: util.GetCategory(170, util.CarbonMonoxide),
		}},
	}

	for _, v := range objects {
		dbb.Add(&v.data)
	}

	dr, err := dbb.Last24H(util.CarbonMonoxide)
	if err != nil {
		t.Fatalf("Expected nil with Last24H, got %v", err)
	}

	if len(*dr) != len(objects) {
		t.Fatalf("Expected length of %d, got %d", len(objects), len(*dr))
	}

	for k, v := range *dr {
		if !v.Equals(&objects[k].expected) {
			t.Error("Objects are not the same")
			t.Errorf("Expected: %v", objects[k].expected)
			t.Errorf("Result: %v", v)
		}
	}
}

// Test DB.Median.
func TestDB_Median(t *testing.T) {
	_ = env.Load(util.EnvTestFilePath)

	dbb := db.DB{}
	dbb.Load()
	_ = dbb.Connect()
	defer dbb.Close()

	_ = dbb.Init()

	creationTime := time.Now().Round(0)

	var objects = [5]model.Data{
		{
			DataType:  util.AirQuality,
			Value:     140,
			Timestamp: creationTime.Add(time.Minute * -2),
		},
		{
			DataType:  util.AirQuality,
			Value:     205,
			Timestamp: creationTime.Add(time.Second * -2),
		},
		{
			DataType:  util.AirQuality,
			Value:     270,
			Timestamp: creationTime.Add(time.Second * -5),
		},
		{
			DataType:  util.AirQuality,
			Value:     33,
			Timestamp: creationTime.Add(time.Second * -21),
		},
		{
			DataType:  util.AirQuality,
			Value:     195,
			Timestamp: creationTime.Add(time.Minute * -1),
		},
	}

	for _, v := range objects {
		dbb.Add(&v)
	}

	d, err := dbb.Median(util.AirQuality)
	if err != nil {
		t.Fatalf("Expected nil with Median, got %v", err)
	}

	dr := model.DataResponse{
		Data: model.Data{
			DataType:  util.AirQuality,
			Value:     195,
			Timestamp: creationTime.Add(time.Minute * -1),
		},
		Category: util.GetCategory(195, util.AirQuality),
	}

	if !d.Equals(&dr) {
		t.Error("Objects are not the same")
		t.Errorf("Expected: %v", dr)
		t.Errorf("Result: %v", d)
	}
}

// Test DB.Max.
func TestDB_Max(t *testing.T) {
	_ = env.Load(util.EnvTestFilePath)

	dbb := db.DB{}
	dbb.Load()
	_ = dbb.Connect()
	defer dbb.Close()

	_ = dbb.Init()

	creationTime := time.Now().Round(0)

	var objects = [5]model.Data{
		{
			DataType:  util.Raindrops,
			Value:     230,
			Timestamp: creationTime.Add(time.Minute * -1),
		},
		{
			DataType:  util.Raindrops,
			Value:     420,
			Timestamp: creationTime.Add(time.Second * -3),
		},
		{
			DataType:  util.Raindrops,
			Value:     114,
			Timestamp: creationTime.Add(time.Second * -7),
		},
		{
			DataType:  util.Raindrops,
			Value:     47,
			Timestamp: creationTime.Add(time.Second * -41),
		},
		{
			DataType:  util.Raindrops,
			Value:     842,
			Timestamp: creationTime.Add(time.Minute * -2),
		},
	}

	for _, v := range objects {
		dbb.Add(&v)
	}

	d, err := dbb.Max(util.Raindrops)
	if err != nil {
		t.Fatalf("Expected nil with Max, got %v", err)
	}

	dr := model.DataResponse{
		Data: model.Data{
			DataType:  util.Raindrops,
			Value:     842,
			Timestamp: creationTime.Add(time.Minute * -2),
		},
		Category: util.GetCategory(842, util.Raindrops),
	}

	if !d.Equals(&dr) {
		t.Error("Objects are not the same")
		t.Errorf("Expected: %v", dr)
		t.Errorf("Result: %v", d)
	}
}

// Test DB.Min.
func TestDB_Min(t *testing.T) {
	_ = env.Load(util.EnvTestFilePath)

	dbb := db.DB{}
	dbb.Load()
	_ = dbb.Connect()
	defer dbb.Close()

	_ = dbb.Init()

	creationTime := time.Now().Round(0)

	var objects = [5]model.Data{
		{
			DataType:  util.SoilMoisture,
			Value:     250,
			Timestamp: creationTime.Add(time.Second * -10),
		},
		{
			DataType:  util.SoilMoisture,
			Value:     307,
			Timestamp: creationTime.Add(time.Second * -17),
		},
		{
			DataType:  util.SoilMoisture,
			Value:     348,
			Timestamp: creationTime.Add(time.Second * -2),
		},
		{
			DataType:  util.SoilMoisture,
			Value:     412,
			Timestamp: creationTime.Add(time.Minute * -5),
		},
		{
			DataType:  util.SoilMoisture,
			Value:     440,
			Timestamp: creationTime.Add(time.Minute * -1),
		},
	}

	for _, v := range objects {
		dbb.Add(&v)
	}

	d, err := dbb.Min(util.SoilMoisture)
	if err != nil {
		t.Fatalf("Expected nil with Min, got %v", err)
	}

	dr := model.DataResponse{
		Data: model.Data{
			DataType:  util.SoilMoisture,
			Value:     250,
			Timestamp: creationTime.Add(time.Second * -10),
		},
		Category: util.GetCategory(250, util.SoilMoisture),
	}

	if !d.Equals(&dr) {
		t.Error("Objects are not the same")
		t.Errorf("Expected: %v", dr)
		t.Errorf("Result: %v", d)
	}
}
