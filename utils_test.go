package main

import (
	"errors"
	"testing"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "Valid HTTP URL",
			url:     "http://example.com",
			wantErr: false,
		},
		{
			name:    "Valid HTTPS URL",
			url:     "https://example.com/path?query=value",
			wantErr: false,
		},
		{
			name:    "Invalid URL - No Scheme",
			url:     "example.com",
			wantErr: true,
		},
		{
			name:    "Invalid URL - Empty",
			url:     "",
			wantErr: true,
		},
		{
			name:    "Invalid URL - Malformed",
			url:     "http://",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLogErr(t *testing.T) {
	// Test with nil error
	LogErr(nil) // Should not panic

	// Test with error
	LogErr(errors.New("test error")) // Should not panic and will print to stdout
}

func TestCheckErr(t *testing.T) {
	// Test with nil error
	CheckErr(nil) // Should not panic

	// Test with error
	defer func() {
		if r := recover(); r == nil {
			t.Error("CheckErr() should have panicked with non-nil error")
		}
	}()
	CheckErr(errors.New("test error")) // Should panic
}
