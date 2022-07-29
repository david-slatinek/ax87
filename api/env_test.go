package main

import (
	"testing"
)

func TestLoad(t *testing.T) {
	err := Load("test.env")
	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}
}

func TestValidate(t *testing.T) {
	err := Validate()
	if err != nil {
		t.Errorf("Expected %v, got %v", nil, err)
	}
}
