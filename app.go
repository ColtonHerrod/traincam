package main

import (
	"archive/zip"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Camera represents our train camera data
type Camera struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	YoutubeURL  string  `json:"youtubeUrl"`
	Description string  `json:"description"`
	Country     string  `json:"country"`
	State       string  `json:"state"`
}

// App struct
type App struct {
	ctx     context.Context
	cameras []Camera
	nextID  int
}

// KML XML structures
type KML struct {
	XMLName    xml.Name    `xml:"http://www.opengis.net/kml/2.2 kml"`
	Document   *Document   `xml:"http://www.opengis.net/kml/2.2 Document"`
	Placemarks []Placemark `xml:"http://www.opengis.net/kml/2.2 Placemark"`
}

type Document struct {
	XMLName    xml.Name    `xml:"http://www.opengis.net/kml/2.2 Document"`
	Placemarks []Placemark `xml:"http://www.opengis.net/kml/2.2 Placemark"`
	Folders    []Folder    `xml:"http://www.opengis.net/kml/2.2 Folder"`
}

type Folder struct {
	XMLName    xml.Name    `xml:"http://www.opengis.net/kml/2.2 Folder"`
	Name       string      `xml:"http://www.opengis.net/kml/2.2 name"`
	Placemarks []Placemark `xml:"http://www.opengis.net/kml/2.2 Placemark"`
}

type Placemark struct {
	XMLName     xml.Name `xml:"http://www.opengis.net/kml/2.2 Placemark"`
	Name        string   `xml:"http://www.opengis.net/kml/2.2 name"`
	Description string   `xml:"http://www.opengis.net/kml/2.2 description"`
	Point       struct {
		XMLName     xml.Name `xml:"http://www.opengis.net/kml/2.2 Point"`
		Coordinates string   `xml:"http://www.opengis.net/kml/2.2 coordinates"`
	} `xml:"http://www.opengis.net/kml/2.2 Point"`
}

// NewApp creates a new App instance
func NewApp() *App {
	app := &App{
		nextID: 1,
		cameras: []Camera{
			{
				ID:         1,
				Name:       "Central Station North",
				Lat:        51.505,
				Lng:        -0.09,
				YoutubeURL: "https://www.youtube.com/embed/dQw4w9WgXcQ",
			},
			{
				ID:         2,
				Name:       "East Junction Crossing",
				Lat:        51.515,
				Lng:        -0.1,
				YoutubeURL: "https://www.youtube.com/embed/dQw4w9WgXcQ",
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
	if err := a.LoadCameras(); err != nil {
		// If loading fails, just keep the hardcoded cameras
		fmt.Printf("Warning: could not load cameras from storage: %v\n", err)
	}
}

// GetCameras returns our list of train cameras to the frontend
func (a *App) GetCameras() []Camera {
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

// ImportKML imports cameras from a KML/KMZ file
func (a *App) ImportKML(filePath string) ([]Camera, error) {
	var kmlContent []byte
	var err error

	if strings.HasSuffix(strings.ToLower(filePath), ".kmz") {
		kmlContent, err = extractKMLFromKMZ(filePath)
	} else {
		kmlContent, err = os.ReadFile(filePath)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Debug: log the KML content
	fmt.Printf("DEBUG: KML content length: %d bytes\n", len(kmlContent))
	fmt.Printf("DEBUG: First 500 chars:\n%s\n", string(kmlContent[:min(500, len(kmlContent))]))

	var kml KML
	err = xml.Unmarshal(kmlContent, &kml)
	if err != nil {
		return nil, fmt.Errorf("failed to parse KML: %w", err)
	}

	// Debug: log what was parsed
	fmt.Printf("DEBUG: KML.Document is nil: %v\n", kml.Document == nil)
	if kml.Document != nil {
		fmt.Printf("DEBUG: Document placemarks: %d\n", len(kml.Document.Placemarks))
		fmt.Printf("DEBUG: Document folders: %d\n", len(kml.Document.Folders))
	}
	fmt.Printf("DEBUG: Root placemarks: %d\n", len(kml.Placemarks))

	var placemarks []Placemark
	if kml.Document != nil && len(kml.Document.Placemarks) > 0 {
		placemarks = kml.Document.Placemarks
	} else if kml.Document != nil && len(kml.Document.Folders) > 0 {
		// Extract placemarks from folders
		for _, folder := range kml.Document.Folders {
			fmt.Printf("DEBUG: Folder '%s': %d placemarks\n", folder.Name, len(folder.Placemarks))
			placemarks = append(placemarks, folder.Placemarks...)
		}
	} else {
		placemarks = kml.Placemarks
	}

	fmt.Printf("DEBUG: Using %d placemarks\n", len(placemarks))

	importedCameras := []Camera{}
	for i, pm := range placemarks {
		fmt.Printf("DEBUG: Placemark %d: name=%s, coords=%s\n", i, pm.Name, pm.Point.Coordinates)

		if pm.Point.Coordinates == "" {
			continue
		}

		lat, lng, err := parseCoordinates(pm.Point.Coordinates)
		if err != nil {
			fmt.Printf("DEBUG: Failed to parse coords for %s: %v\n", pm.Name, err)
			continue
		}

		// If no YouTube URL found in description, use a default placeholder
		youtubeUrl := extractYoutubeURL(pm.Description)
		if youtubeUrl == "" {
			youtubeUrl = "https://www.youtube.com/embed/dQw4w9WgXcQ"
		}

		camera := Camera{

			ID:          a.nextID,
			Name:        pm.Name,
			Lat:         lat,
			Lng:         lng,
			Description: pm.Description,
			YoutubeURL:  youtubeUrl,
		}

		a.cameras = append(a.cameras, camera)
		importedCameras = append(importedCameras, camera)
		a.nextID++
	}

	if len(importedCameras) == 0 {
		return nil, fmt.Errorf("no valid placemarks found in KML file")
	}

	// Save cameras to persistent storage after successful import
	if err := a.SaveCameras(); err != nil {
		fmt.Printf("Warning: could not save cameras to storage: %v\n", err)
	}

	return importedCameras, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// extractKMLFromKMZ extracts the doc.kml file from a KMZ (zipped KML)
func extractKMLFromKMZ(kmzPath string) ([]byte, error) {
	reader, err := zip.OpenReader(kmzPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open KMZ file: %w", err)
	}
	defer reader.Close()

	var foundFiles []string
	for _, file := range reader.File {
		foundFiles = append(foundFiles, file.Name)

		// Look for .kml file (case-insensitive)
		if strings.HasSuffix(strings.ToLower(file.Name), ".kml") {
			f, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to open KML file in archive: %w", err)
			}
			defer f.Close()

			content, err := io.ReadAll(f)
			if err != nil {
				return nil, fmt.Errorf("failed to read KML content: %w", err)
			}
			return content, nil
		}
	}

	return nil, fmt.Errorf("no KML file found in KMZ archive. Found files: %v", foundFiles)
}

// extractYoutubeURL extracts and converts a YouTube URL from a description
// Returns the standard embed URL which works best in Wails
func extractYoutubeURL(description string) string {
	// Method 1: Look for youtube.com URLs with different patterns
	// Handle "http://", "https://", and URLs with trailing punctuation
	urlPatterns := []string{
		"https://www.youtube.com/watch?v=",
		"http://www.youtube.com/watch?v=",
		"youtube.com/watch?v=",
		"https://youtu.be/",
		"http://youtu.be/",
		"youtu.be/",
	}

	for _, pattern := range urlPatterns {
		startIdx := strings.Index(description, pattern)
		if startIdx != -1 {
			// Find the video ID (starts after pattern)
			videoIDStart := startIdx + len(pattern)
			// Find the end of the video ID (until space, quote, comma, period, question mark, &, <, >, /, =, or end of string)
			videoIDEnd := videoIDStart
			for videoIDEnd < len(description) && description[videoIDEnd] != ' ' &&
				description[videoIDEnd] != '"' && description[videoIDEnd] != '"' &&
				description[videoIDEnd] != '\'' && description[videoIDEnd] != ',' &&
				description[videoIDEnd] != '.' && description[videoIDEnd] != '?' &&
				description[videoIDEnd] != '&' && description[videoIDEnd] != '<' &&
				description[videoIDEnd] != '>' && description[videoIDEnd] != '\n' &&
				description[videoIDEnd] != '\r' && description[videoIDEnd] != '/' &&
				description[videoIDEnd] != '=' {
				videoIDEnd++
			}
			videoID := description[videoIDStart:videoIDEnd]
			return "https://www.youtube.com/embed/" + videoID + "?modestbranding=1"
		}
	}

	return ""
}

// parseCoordinates parses KML coordinates format: "lng,lat[,altitude]"
func parseCoordinates(coordStr string) (float64, float64, error) {
	coords := strings.Split(strings.TrimSpace(coordStr), ",")
	if len(coords) < 2 {
		return 0, 0, fmt.Errorf("invalid coordinates format")
	}

	lng, err := strconv.ParseFloat(coords[0], 64)
	if err != nil {
		return 0, 0, err
	}

	lat, err := strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return 0, 0, err
	}

	return lat, lng, nil
}

// getCamerasStoragePath returns the path to the cameras.json file
func getCamerasStoragePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	storageDir := filepath.Join(homeDir, ".traincam")
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(storageDir, "cameras.json"), nil
}

// SaveCameras saves the current camera list to a JSON file
func (a *App) SaveCameras() error {
	storagePath, err := getCamerasStoragePath()
	if err != nil {
		return fmt.Errorf("could not determine storage path: %w", err)
	}

	data, err := json.MarshalIndent(a.cameras, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal cameras to JSON: %w", err)
	}

	if err := os.WriteFile(storagePath, data, 0644); err != nil {
		return fmt.Errorf("could not write cameras file: %w", err)
	}

	return nil
}

// LoadCameras loads the camera list from a JSON file
func (a *App) LoadCameras() error {
	storagePath, err := getCamerasStoragePath()
	if err != nil {
		return fmt.Errorf("could not determine storage path: %w", err)
	}

	// If file doesn't exist, keep the hardcoded cameras
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(storagePath)
	if err != nil {
		return fmt.Errorf("could not read cameras file: %w", err)
	}

	var cameras []Camera
	if err := json.Unmarshal(data, &cameras); err != nil {
		return fmt.Errorf("could not parse cameras JSON: %w", err)
	}

	// Replace in-memory cameras and update nextID
	a.cameras = cameras
	a.nextID = 1
	for _, cam := range a.cameras {
		if cam.ID >= a.nextID {
			a.nextID = cam.ID + 1
		}
	}

	return nil
}
