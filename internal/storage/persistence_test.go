package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetCamerasStoragePath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "traincam_test_paths")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	path, err := GetCamerasStoragePathWithBase(tmpDir)
	if err != nil {
		t.Fatalf("Error getting path: %v", err)
	}

	expectedSuffix := "cameras.json"
	if !strings.HasSuffix(path, expectedSuffix) {
		t.Errorf("Expected path to end with %s, got %s", expectedSuffix, path)
	}

	// Verify directory exists
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("Storage directory does not exist: %v", err)
	}
}

func TestSaveAndLoadCameras(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "traincam_test_data")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	original := []Camera{
		{ID: 1, Name: "Test Camera", Lat: 10.0, Lng: 20.0},
	}

	if err := SaveCameras(original, tmpDir); err != nil {
		t.Fatalf("Failed to save cameras: %v", err)
	}

	loaded, err := LoadCameras(tmpDir)
	if err != nil {
		t.Fatalf("Failed to load cameras: %v", err)
	}

	if len(loaded) != 1 || loaded[0].Name != "Test Camera" {
		t.Errorf("Loaded cameras do not match original. Got: %+v", loaded)
	}
}
