package main

import (
	"testing"
	"time"
)

// Test DB.LoadFields.
func TestDB_LoadFields(t *testing.T) {
	_ = Load("test.env")

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
	_ = Load("test.env")

	db := DB{}
	db.LoadFields()
	defer db.client.Close()

	if err := db.Connect(); err != nil {
		t.Fatalf("Expected nil with Connect, got %v", err)
	}
}

// Test DB.Init
func TestDB_Init(t *testing.T) {
	_ = Load("test.env")

	db := DB{}
	db.LoadFields()
	_ = db.Connect()
	defer db.client.Close()

	if err := db.Init(); err != nil {
		t.Fatalf("Expected nil with Init, got %v", err)
	}
}

// Test MapCO2.
func TestMapCO2(t *testing.T) {
	var tests = []struct {
		value, expected int
	}{
		{-10, 1},
		{0, 1},
		{17, 1},
		{30, 1},
		{42, 2},
		{69, 2},
		{85, 3},
		{120, 3},
		{150, 3},
		{151, 4},
		{188, 4},
		{200, 4},
		{217, 5},
		{300, 5},
		{361, 5},
		{400, 5},
		{470, 6},
		{600, 6},
		{707, 6},
		{747, 6},
		{800, 6},
		{801, 7},
		{900, 7},
		{1000, 7},
	}

	for _, value := range tests {
		ans := MapCO2(value.value)
		if ans != value.expected {
			t.Errorf("Expected %d with MapCO2(%d), got %d", value.expected, value.value, ans)
		}
	}
}

// Test MapAir.
func TestMapAir(t *testing.T) {
	var tests = []struct {
		value, expected int
	}{
		{-10, 1},
		{0, 1},
		{17, 1},
		{30, 1},
		{50, 1},
		{69, 2},
		{85, 2},
		{100, 2},
		{120, 3},
		{150, 3},
		{188, 4},
		{200, 4},
		{217, 5},
		{300, 5},
		{361, 6},
		{400, 6},
	}

	for _, value := range tests {
		ans := MapAir(value.value)
		if ans != value.expected {
			t.Errorf("Expected %d with MapAir(%d), got %d", value.expected, value.value, ans)
		}
	}
}

// Test MapValue.
func TestMapValue(t *testing.T) {
	var tests = []struct {
		x, inMin, inMax, outMin, outMax float64
		expected                        int
	}{
		{0, 0, 1024, 1, 4, 1},
		{1024, 0, 1024, 1, 4, 4},
		{10, 0, 1024, 1, 4, 1},
		{123, 0, 1024, 1, 4, 1},
		{270, 0, 1024, 1, 4, 2},
		{350, 0, 1024, 1, 4, 2},
		{512, 0, 1024, 1, 4, 3},
		{645, 0, 1024, 1, 4, 3},
		{749, 0, 1024, 1, 4, 3},
		{750, 0, 1024, 1, 4, 3},
		{751, 0, 1024, 1, 4, 3},
		{891, 0, 1024, 1, 4, 4},
		{955, 0, 1024, 1, 4, 4},
		{1000, 0, 1024, 1, 4, 4},
		{1015, 0, 1024, 1, 4, 4},
		{238, 489, 238, 0, 100, 100},
		{489, 489, 238, 0, 100, 0},
		{300, 489, 238, 0, 100, 75},
		{250, 489, 238, 0, 100, 95},
		{387, 489, 238, 0, 100, 41},
		{480, 489, 238, 0, 100, 4},
		{400, 489, 238, 0, 100, 35},
		{290, 489, 238, 0, 100, 79},
	}

	for _, value := range tests {
		ans := MapValue(value.x, value.inMin, value.inMax, value.outMin, value.outMax)
		if ans != value.expected {
			t.Errorf("Expected %d with MapValue(%f), got %d", value.expected, value.x, ans)
		}
	}
}

// Test DB.Latest.
func TestDB_Latest(t *testing.T) {
	_ = Load("test.env")

	db := DB{}
	db.LoadFields()
	_ = db.Connect()
	defer db.client.Close()

	_ = db.Init()

	creationTime := time.Now().Round(0)

	var objects = []struct {
		data     Data
		expected DataResponse
	}{
		{Data{
			DataType:  carbonMonoxide,
			Value:     45,
			TimeStamp: creationTime,
		}, DataResponse{
			Data: Data{
				DataType:  carbonMonoxide,
				Value:     45,
				TimeStamp: creationTime,
			},
			Category: MapCO2(45),
		}},
		{Data{
			DataType:  airQuality,
			Value:     125,
			TimeStamp: creationTime,
		}, DataResponse{
			Data: Data{
				DataType:  airQuality,
				Value:     125,
				TimeStamp: creationTime,
			},
			Category: MapAir(125),
		}},
		{Data{
			DataType:  raindrops,
			Value:     800,
			TimeStamp: creationTime,
		}, DataResponse{
			Data: Data{
				DataType:  raindrops,
				Value:     800,
				TimeStamp: creationTime,
			},
			Category: MapValue(800, 0, 1024, 1, 4),
		}},
		{Data{
			DataType:  soilMoisture,
			Value:     400,
			TimeStamp: creationTime,
		}, DataResponse{
			Data: Data{
				DataType:  soilMoisture,
				Value:     400,
				TimeStamp: creationTime,
			},
			Category: MapValue(400, 489, 238, 0, 100),
		}},
	}

	for _, v := range objects {
		db.Add(&v.data)
	}

	types := [4]string{carbonMonoxide, airQuality, raindrops, soilMoisture}

	for k, v := range types {
		dr, err := db.Latest(v)
		if err != nil {
			t.Fatalf("Expected nil with Latest, got %v", err)
		}

		if !dr.Compare(&objects[k].expected) {
			t.Error("Objects not the same")
			t.Errorf("Expected: %v", objects[k].expected)
			t.Errorf("Result: %v", dr)
		}
	}

	_ = db.Init()
}

// Test DB.Last24H.
func TestDB_Last24H(t *testing.T) {
	_ = Load("test.env")

	db := DB{}
	db.LoadFields()
	_ = db.Connect()
	defer db.client.Close()

	_ = db.Init()

	creationTime := time.Now().Round(0)

	var objects = []struct {
		data     Data
		expected DataResponse
	}{
		{Data{
			DataType:  carbonMonoxide,
			Value:     250,
			TimeStamp: creationTime,
		}, DataResponse{
			Data: Data{
				DataType:  carbonMonoxide,
				Value:     250,
				TimeStamp: creationTime,
			},
			Category: MapCO2(250),
		}},
		{Data{
			DataType:  carbonMonoxide,
			Value:     55,
			TimeStamp: creationTime.Add(time.Second * -1),
		}, DataResponse{
			Data: Data{
				DataType:  carbonMonoxide,
				Value:     55,
				TimeStamp: creationTime.Add(time.Second * -1),
			},
			Category: MapCO2(55),
		}},
		{Data{
			DataType:  carbonMonoxide,
			Value:     420,
			TimeStamp: creationTime.Add(time.Second * -10),
		}, DataResponse{
			Data: Data{
				DataType:  carbonMonoxide,
				Value:     420,
				TimeStamp: creationTime.Add(time.Second * -10),
			},
			Category: MapCO2(420),
		}},
		{Data{
			DataType:  carbonMonoxide,
			Value:     69,
			TimeStamp: creationTime.Add(time.Minute * -1),
		}, DataResponse{
			Data: Data{
				DataType:  carbonMonoxide,
				Value:     69,
				TimeStamp: creationTime.Add(time.Minute * -1),
			},
			Category: MapCO2(69),
		}},
		{Data{
			DataType:  carbonMonoxide,
			Value:     170,
			TimeStamp: creationTime.Add(time.Minute * -2),
		}, DataResponse{
			Data: Data{
				DataType:  carbonMonoxide,
				Value:     170,
				TimeStamp: creationTime.Add(time.Minute * -2),
			},
			Category: MapCO2(170),
		}},
	}

	for _, v := range objects {
		db.Add(&v.data)
	}

	dr, err := db.Last24H(carbonMonoxide)
	if err != nil {
		t.Fatalf("Expected nil with Last24H, got %v", err)
	}

	if len(*dr) != len(objects) {
		t.Fatalf("Expected length %d, got %d", len(objects), len(*dr))
	}

	for k, v := range *dr {
		if !v.Compare(&objects[k].expected) {
			t.Error("Objects not the same")
			t.Errorf("Expected: %v", objects[k].expected)
			t.Errorf("Result: %v", v)
		}
	}

	_ = db.Init()
}
