package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Camera represents our train camera data
type Camera struct {
	ID                   int     `json:"id"`
	Name                 string  `json:"name"`
	Lat                  float64 `json:"lat"`
	Lng                  float64 `json:"lng"`
	YoutubeURL           string  `json:"youtubeUrl"`
	Description          string  `json:"description"`
	Country              string  `json:"country"`
	State                string  `json:"state"`
	IconURL              string  `json:"iconUrl"`
	SubscriptionRequired bool    `json:"subscriptionRequired"`
}

// GetCamerasStoragePath returns the path to the cameras.json file in the user's home directory
func GetCamerasStoragePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return GetCamerasStoragePathWithBase(homeDir)
}

// GetCamerasStoragePathWithBase returns the path to cameras.json relative to a base directory
func GetCamerasStoragePathWithBase(base string) (string, error) {
	storageDir := filepath.Join(base, ".traincam")
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(storageDir, "cameras.json"), nil
}

// SaveCameras saves the current camera list to a JSON file. Pass base path for testing.
func SaveCameras(cameras []Camera, base string) error {
	var storagePath string
	var err error

	if base != "" {
		storagePath, err = GetCamerasStoragePathWithBase(base)
	} else {
		storagePath, err = GetCamerasStoragePath()
	}
	
	if err != nil {
		return fmt.Errorf("could not determine storage path: %w", err)
	}

	data, err := json.MarshalIndent(cameras, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal cameras to JSON: %w", err)
	}

	if err := os.WriteFile(storagePath, data, 0644); err != nil {
		return fmt.Errorf("could not write cameras file: %w", err)
	}

	return nil
}

// LoadCameras loads the camera list from a JSON file. Pass base path for testing.
func LoadCameras(base string) ([]Camera, error) {
	var storagePath string
	var err error

	if base != "" {
		storagePath, err = GetCamerasStoragePathWithBase(base)
	} else {
		storagePath, err = GetCamerasStoragePath()
	}
    
	if err != nil {
		return nil, fmt.Errorf("could not determine storage path: %w", err)
	}

	// If file doesn't exist, return nil (let caller handle default)
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		return nil, nil
	}

	data, err := os.ReadFile(storagePath)
	if err != nil {
		return nil, fmt.Errorf("could not read cameras file: %w", err)
	}

	var cameras []Camera
	if err := json.Unmarshal(data, &cameras); err != nil {
		return nil, fmt.Errorf("could not parse cameras JSON: %w", err)
	}

	return cameras, nil
}
