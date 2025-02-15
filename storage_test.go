package main

import (
	"os"
	"testing"
)

func TestStorage(t *testing.T) {
	// Setup test environment
	tempDir, err := os.MkdirTemp("", "gry-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Override storage path for testing
	originalPath := STORAGE_PATH
	STORAGE_PATH = tempDir
	defer func() { STORAGE_PATH = originalPath }()

	// Test CreateURL
	t.Run("CreateURL", func(t *testing.T) {
		// Test successful creation
		err := CreateURL("test1", "https://example.com")
		if err != nil {
			t.Errorf("CreateURL failed: %v", err)
		}

		// Test duplicate creation
		err = CreateURL("test1", "https://example.com")
		if err == nil || err.Error() != "slug already exists" {
			t.Errorf("Expected 'slug already exists' error, got: %v", err)
		}
	})

	// Test GetURL
	t.Run("GetURL", func(t *testing.T) {
		// Test successful retrieval
		url, err := GetURL("test1")
		if err != nil {
			t.Errorf("GetURL failed: %v", err)
		}
		if url != "https://example.com" {
			t.Errorf("Expected 'https://example.com', got '%s'", url)
		}

		// Test non-existent slug
		_, err = GetURL("nonexistent")
		if err == nil {
			t.Error("Expected error for non-existent slug, got nil")
		}
	})

	// Test UpdateURL
	t.Run("UpdateURL", func(t *testing.T) {
		// Test successful update
		err := UpdateURL("test1", "https://updated.com")
		if err != nil {
			t.Errorf("UpdateURL failed: %v", err)
		}

		// Verify update
		url, _ := GetURL("test1")
		if url != "https://updated.com" {
			t.Errorf("Expected 'https://updated.com', got '%s'", url)
		}

		// Test update on non-existent slug
		err = UpdateURL("nonexistent", "https://example.com")
		if err == nil || err.Error() != "slug does not exist" {
			t.Errorf("Expected 'slug does not exist' error, got: %v", err)
		}
	})

	// Test DeleteURL
	t.Run("DeleteURL", func(t *testing.T) {
		// Test successful deletion
		err := DeleteURL("test1")
		if err != nil {
			t.Errorf("DeleteURL failed: %v", err)
		}

		// Verify deletion
		_, err = GetURL("test1")
		if err == nil {
			t.Error("Expected error after deletion, got nil")
		}

		// Test delete non-existent slug
		err = DeleteURL("nonexistent")
		if err == nil || err.Error() != "slug does not exist" {
			t.Errorf("Expected 'slug does not exist' error, got: %v", err)
		}
	})
}
