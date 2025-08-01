package overlay

import (
	"testing"
	"time"

	"github.com/bluenviron/mediamtx/internal/gps"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	if config.Enabled != false {
		t.Errorf("Expected Enabled to be false, got %v", config.Enabled)
	}
	
	if config.DatabaseHost != "localhost" {
		t.Errorf("Expected DatabaseHost to be localhost, got %v", config.DatabaseHost)
	}
	
	if config.UpdateInterval != time.Second {
		t.Errorf("Expected UpdateInterval to be 1s, got %v", config.UpdateInterval)
	}
}

// MockGPSDataProvider is a mock implementation of gps.DataProvider for testing.
type MockGPSDataProvider struct{}

func (m *MockGPSDataProvider) GetCurrentGPS() *gps.Data {
	return &gps.Data{
		Timestamp: time.Now(),
		Latitude:  12.345678,
		Longitude: 123.456789,
		Status:    "A",
	}
}

func TestNewEngineDisabled(t *testing.T) {
	config := DefaultConfig()
	config.Enabled = false
	
	mockProvider := &MockGPSDataProvider{}
	engine, err := NewEngine(config, mockProvider)
	if err != nil {
		t.Fatalf("Failed to create disabled overlay engine: %v", err)
	}
	defer engine.Close()
	
	if engine.IsEnabled() {
		t.Error("Expected overlay engine to be disabled")
	}
}

func TestFormatCoordinate(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{12.345678, "12.345678"},
		{123.456789, "123.456789"},
		{-45.123456, "-45.123456"},
		{0.0, "0.000000"},
	}
	
	for _, test := range tests {
		result := gps.FormatCoordinate(test.input)
		if result != test.expected {
			t.Errorf("gps.FormatCoordinate(%f) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestGPSDataCreation(t *testing.T) {
	gpsData := &gps.Data{
		Timestamp: time.Now(),
		Latitude:  12.345678,
		Longitude: 123.456789,
		Status:    "A",
	}
	
	if gpsData.Latitude != 12.345678 {
		t.Errorf("Expected Latitude to be 12.345678, got %f", gpsData.Latitude)
	}
	
	if gpsData.Status != "A" {
		t.Errorf("Expected Status to be A, got %s", gpsData.Status)
	}
}