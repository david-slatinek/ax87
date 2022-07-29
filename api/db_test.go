package main

import (
	"testing"
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

	if err := db.Connect(); err != nil {
		t.Fatalf("Expected nil with Connect, got %v", err)
	}

	db.client.Close()
}

// Test DB.Init
func TestDB_Init(t *testing.T) {
	_ = Load("test.env")

	db := DB{}
	db.LoadFields()
	_ = db.Connect()

	if err := db.Init(); err != nil {
		t.Fatalf("Expected nil with Init, got %v", err)
	}

	db.client.Close()
}
