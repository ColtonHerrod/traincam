package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"os"
	"regexp"
	"strings"

	"traincam/internal/kml"
	"traincam/internal/storage"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx     context.Context
	cameras []storage.Camera
	nextID  int
}

// NewApp creates a new App instance
func NewApp() *App {
	app := &App{
		nextID: 1,
		cameras: []storage.Camera{
			{
				ID:                   1,
				Name:                 "Central Station North",
				Lat:                  51.505,
				Lng:                  -0.09,
				YoutubeURL:           "https://www.youtube.com/embed/dQw4w9WgXcQ",
				SubscriptionRequired: false,
			},
			{
				ID:                   2,
				Name:                 "East Junction Crossing",
				Lat:                  51.515,
				Lng:                  -0.1,
				YoutubeURL:           "https://www.youtube.com/embed/dQw4w9WgXcQ",
				SubscriptionRequired: false,
			},
		},
	}
	app.nextID = 3
	return app
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// Load cameras from persistent storage
	if loaded, err := storage.LoadCameras(""); err != nil {
		fmt.Printf("Warning: could not load cameras from storage: %v\n", err)
	} else if loaded != nil {
		a.cameras = loaded
		a.nextID = 1
		for _, cam := range a.cameras {
			if cam.ID >= a.nextID {
				a.nextID = cam.ID + 1
			}
		}
	}
}

// GetCameras returns our list of train cameras to the frontend
func (a *App) GetCameras() []storage.Camera {
	return a.cameras
}

// OpenURL opens a URL in the default browser
func (a *App) OpenURL(url string) {
	if a.ctx == nil {
		return
	}
	runtime.BrowserOpenURL(a.ctx, url)
}

// SelectKMLFile opens a native file dialog to select a KML/KMZ file
func (a *App) SelectKMLFile() (string, error) {
	if a.ctx == nil {
		return "", fmt.Errorf("application context not initialized")
	}

	options := runtime.OpenDialogOptions{
		Title: "Select KML/KMZ File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "KML/KMZ Files (*.kml, *.kmz)",
				Pattern:     "*.kml;*.kmz",
			},
			{
				DisplayName: "All Files (*.*)",
				Pattern:     "*",
			},
		},
	}
	path, err := runtime.OpenFileDialog(a.ctx, options)
	if err != nil {
		return "", fmt.Errorf("failed to open file dialog: %w", err)
	}
	return path, nil
}

func extractState(name string) string {
	// Matches ", XX" or ", XXXX" (state name) at the end of a string, allowing for optional trailing info or whitespace.
	re := regexp.MustCompile(`,\s*([A-Z]{2}|[A-Z][a-z]+)(?:\s*\(.*\))?\s*$`)
	matches := re.FindStringSubmatch(name)
	if len(matches) > 1 {
		state := matches[1]
		// If it's a full name, map it to the shortcode
		if !isShortcode(state) {
			if shortcode, ok := stateToShortcode[state]; ok {
				return shortcode
			}
		}
		return state
	}
	return ""
}

var stateToShortcode = map[string]string{
	"Alabama": "AL", "Alaska": "AK", "Arizona": "AZ", "Arkansas": "AR", "California": "CA",
	"Colorado": "CO", "Connecticut": "CT", "Delaware": "DE", "Florida": "FL", "Georgia": "GA",
	"Hawaii": "HI", "Idaho": "ID", "Illinois": "IL", "Indiana": "IN", "Iowa": "IA",
	"Kansas": "KS", "Kentucky": "KY", "Louisiana": "LA", "Maine": "ME", "Maryland": "MD",
	"Massachusetts": "MA", "Michigan": "MI", "Minnesota": "MN", "Mississippi": "MS", "Missouri": "MO",
	"Montana": "MT", "Nebraska": "NE", "Nevada": "NV", "New Hampshire": "NH", "New Jersey": "NJ",
	"New Mexico": "NM", "New York": "NY", "North Carolina": "NC", "North Dakota": "ND", "Ohio": "OH",
	"Oklahoma": "OK", "Oregon": "OR", "Pennsylvania": "PA", "Rhode Island": "RI", "South Carolina": "SC",
	"South Dakota": "SD", "Tennessee": "TN", "Texas": "TX", "Utah": "UT", "Vermont": "VT",
	"Virginia": "VA", "Washington": "WA", "West Virginia": "WV", "Wisconsin": "WI", "Wyoming": "WY",
}

func isShortcode(s string) bool {
	return len(s) == 2 && s[0] >= 'A' && s[0] <= 'Z' && s[1] >= 'A' && s[1] <= 'Z'
}

// ImportKML imports cameras from a KML/KMZ file
func (a *App) ImportKML(filePath string) ([]storage.Camera, error) {
	var kmlContent []byte
	var filesMap map[string][]byte
	var err error

	if strings.HasSuffix(strings.ToLower(filePath), ".kmz") {
		fmt.Printf("DEBUG: Extracting KMZ file: %s\n", filePath)
		kmlContent, filesMap, err = kml.ExtractKMLFromKMZ(filePath)
		fmt.Printf("DEBUG: Extracted KML content size: %d, Files in map: %d\n", len(kmlContent), len(filesMap))
	} else {

		kmlContent, err = os.ReadFile(filePath)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var kmlData kml.KML
	if err := xml.Unmarshal(kmlContent, &kmlData); err != nil {
		return nil, fmt.Errorf("failed to parse KML: %w", err)
	}

	var placemarks []kml.Placemark
	if kmlData.Document != nil && len(kmlData.Document.Placemarks) > 0 {
		placemarks = kmlData.Document.Placemarks
	} else if kmlData.Document != nil && len(kmlData.Document.Folders) > 0 {
		for _, folder := range kmlData.Document.Folders {
			placemarks = append(placemarks, folder.Placemarks...)
		}
	} else {
		placemarks = kmlData.Placemarks
	}

	importedCameras := []storage.Camera{}
	for _, pm := range placemarks {
		if pm.Point.Coordinates == "" {
			continue
		}

		lat, lng, err := kml.ParseCoordinates(pm.Point.Coordinates)
		if err != nil {
			continue
		}

		youtubeUrl := kml.ExtractYoutubeURL(pm.Description)
		if youtubeUrl == "" {
			youtubeUrl = "https://www.youtube.com/embed/dQw4w9WgXcQ"
		}

		camera := storage.Camera{
			ID:                   a.nextID,
			Name:                 pm.Name,
			Lat:                  lat,
			Lng:                  lng,
			Description:          pm.Description,
			YoutubeURL:           youtubeUrl,
			SubscriptionRequired: false,
			State:                extractState(pm.Name),
		}

		a.cameras = append(a.cameras, camera)
		importedCameras = append(importedCameras, camera)
		a.nextID++
	}

	if len(importedCameras) == 0 {
		return nil, fmt.Errorf("no valid placemarks found in KML file")
	}

	if err := storage.SaveCameras(a.cameras, ""); err != nil {
		fmt.Printf("Warning: could not save cameras to storage: %v\n", err)
	}

	return importedCameras, nil
}

// RemoveCamera removes a camera from the list and persistent storage by ID
func (a *App) RemoveCamera(id int) []storage.Camera {
	var updatedCameras []storage.Camera
	for _, cam := range a.cameras {
		if cam.ID != id {
			updatedCameras = append(updatedCameras, cam)
		}
	}

	a.cameras = updatedCameras

	// Update the persistent storage
	if err := storage.SaveCameras(a.cameras, ""); err != nil {
		fmt.Printf("Warning: could not save cameras after removal: %v\n", err)
	}

	return a.cameras
}
