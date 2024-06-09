package main

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	// Test valid time format (HH:MM)
	validTime := "12:34"
	expectedTime := time.Date(0, 1, 1, 12, 34, 0, 0, time.UTC)
	parsedTime, err := parseTime(validTime)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !parsedTime.Equal(expectedTime) {
		t.Errorf("Expected time %v, but got %v", expectedTime, parsedTime)
	}

	// Test valid time format (HH:MM:SS)
	validTimeWithSeconds := "12:34:56"
	expectedTimeWithSeconds := time.Date(0, 1, 1, 12, 34, 56, 0, time.UTC)
	parsedTimeWithSeconds, err := parseTime(validTimeWithSeconds)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !parsedTimeWithSeconds.Equal(expectedTimeWithSeconds) {
		t.Errorf("Expected time %v, but got %v", expectedTimeWithSeconds, parsedTimeWithSeconds)
	}

	// Test invalid time format
	invalidTime := "12:34:56:78"
	_, err = parseTime(invalidTime)
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}
}

func TestGetVersion(t *testing.T) {
	expectedVersionByte, err := os.ReadFile("VERSION")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expectedVersion := strings.TrimSpace(string(expectedVersionByte))

	version := GetVersion()
	if version != expectedVersion {
		t.Errorf("Expected version %s, but got %s", expectedVersion, version)
	}
}
