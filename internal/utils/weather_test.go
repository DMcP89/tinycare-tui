package utils

import "testing"

func TestValidZipCode(t *testing.T) {
	_, err := GetWeather("07070")
	if err != nil {
		t.Error(err)
	}
}

func TestEmptyStringZipCode(t *testing.T) {
	_, err := GetWeather("")
	if err == nil {
		t.Fatalf("Expected error for empty string zip")
	}
}

func TestInvalidZipCode(t *testing.T) {
	_, err := GetWeather("ABCDEF")
	if err == nil {
		t.Fatalf("Expected error for invalid zip")
	}
}
