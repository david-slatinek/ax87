package util_test

import (
	"api/util"
	"testing"
)

// Test util.MapCO2.
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
		ans := util.MapCO2(value.value)
		if ans != value.expected {
			t.Errorf("Expected %d with util.MapCO2(%d), got %d", value.expected, value.value, ans)
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
		ans := util.MapAir(value.value)
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
		ans := util.MapValue(value.x, value.inMin, value.inMax, value.outMin, value.outMax)
		if ans != value.expected {
			t.Errorf("Expected %d with MapValue(%f), got %d", value.expected, value.x, ans)
		}
	}
}

// Test GetCategory
func TestGetCategory(t *testing.T) {
	var tests = []struct {
		value, expected int
		dataType        string
	}{
		{250, 5, util.CarbonMonoxide},
		{80, 3, util.CarbonMonoxide},
		{434, 6, util.CarbonMonoxide},
		{45, 1, util.AirQuality},
		{220, 5, util.AirQuality},
		{197, 4, util.AirQuality},
		{305, 2, util.Raindrops},
		{921, 4, util.Raindrops},
		{817, 3, util.Raindrops},
		{280, 83, util.SoilMoisture},
		{330, 63, util.SoilMoisture},
		{451, 15, util.SoilMoisture},
	}

	for _, value := range tests {
		ans := util.GetCategory(value.value, value.dataType)
		if ans != value.expected {
			t.Errorf("Expected %d with GetCategory(%d), got %d", value.expected, value.value, ans)
		}
	}
}
