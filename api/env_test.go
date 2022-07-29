package main

import (
	"testing"
)

// Test Load.
func TestLoad(t *testing.T) {
	err := Load("test.env")
	if err != nil {
		t.Errorf("Expected nil with Load, got %v", err)
	}
}

// Test Validate.
func TestValidate(t *testing.T) {
	err := Validate()
	if err != nil {
		t.Errorf("Expected nil with Validate, got %v", err)
	}
}
