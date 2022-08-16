package cache_test

import (
	"api/cache"
	"api/env"
	"api/model"
	"api/util"
	"encoding/json"
	"testing"
	"time"
)

// Test Cache.Create.
func TestCache_Create(t *testing.T) {
	_ = env.Load(util.EnvTestFilePath)

	ca := &cache.Cache{}
	ca.Load()
	defer func(ca *cache.Cache) {
		_ = ca.Close()
	}(ca)

	if err := ca.Create(); err != nil {
		t.Errorf("Expected nil with Create, got %v", err)
	}
}

// Test Cache.Close.
func TestCache_Close(t *testing.T) {
	_ = env.Load(util.EnvTestFilePath)

	ca := cache.Cache{}
	ca.Load()

	_ = ca.Create()

	if err := ca.Close(); err != nil {
		t.Errorf("Expected nil with Close, got %v", err)
	}
}

// Test Cache.Add.
func TestCache_Add(t *testing.T) {
	_ = env.Load(util.EnvTestFilePath)

	ca := &cache.Cache{}
	ca.Load()
	defer func(ca *cache.Cache) {
		_ = ca.Close()
	}(ca)

	_ = ca.Create()

	dr := &model.DataResponse{
		Data: model.Data{
			DataType:  util.SoilMoisture,
			Value:     305,
			Timestamp: time.Now(),
		},
		Category: util.GetCategory(305, util.SoilMoisture),
	}

	if err := ca.Add(dr); err != nil {
		t.Errorf("Expected nil with Add, got %v", err)
	}
}

// Test Cache.Get.
func TestCache_Get(t *testing.T) {
	_ = env.Load(util.EnvTestFilePath)

	ca := &cache.Cache{}
	ca.Load()
	defer func(ca *cache.Cache) {
		_ = ca.Close()
	}(ca)

	_ = ca.Create()

	dr := &model.DataResponse{
		Data: model.Data{
			DataType:  util.Raindrops,
			Value:     100,
			Timestamp: time.Now(),
		},
		Category: util.GetCategory(100, util.Raindrops),
	}
	_ = ca.Add(dr)

	value, err := ca.Get(util.Raindrops)

	if err != nil {
		t.Fatalf("Expected nil with Get, got %v", err)
	}

	var data *model.DataResponse

	_ = json.Unmarshal([]byte(value), &data)

	if !dr.Equals(data) {
		t.Error("Objects are not the same")
		t.Errorf("Expected: %v", dr)
		t.Errorf("Result: %v", data)
	}
}
