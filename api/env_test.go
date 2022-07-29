package main

import (
	"testing"
)

// Test Load.
func TestLoad(t *testing.T) {
	if err := Load("test.env"); err != nil {
		t.Errorf("Expected nil with Load, got %v", err)
	}
}

// Test Validate.
func TestValidate(t *testing.T) {
	if err := Validate(); err != nil {
		t.Errorf("Expected nil with Validate, got %v", err)
	}
}
