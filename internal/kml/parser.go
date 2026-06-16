package kml

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
)

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

// ExtractKMLFromKMZ extracts the doc.kml file from a KMZ (zipped KML)
func ExtractKMLFromKMZ(kmzPath string) ([]byte, map[string][]byte, error) {

	reader, err := zip.OpenReader(kmzPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open KMZ file: %w", err)
	}
	defer reader.Close()

	var foundFiles []string
	filesMap := make(map[string][]byte)
	var kmlContent []byte

	for _, file := range reader.File {
		foundFiles = append(foundFiles, file.Name)
		f, err := file.Open()
		if err != nil {
			continue
		}
		content, _ := io.ReadAll(f)
		f.Close()

		if strings.HasSuffix(strings.ToLower(file.Name), ".kml") {
			fmt.Printf("DEBUG: Found KML file: %s\n", file.Name)
			kmlContent = content
		} else {
			fmt.Printf("DEBUG: Found other file: %s (size: %d)\n", file.Name, len(content))
			filesMap[file.Name] = content
		}
	}

	if kmlContent == nil {
		return nil, nil, fmt.Errorf("no KML file found in KMZ archive. Found files: %v", foundFiles)
	}

	fmt.Printf("DEBUG: Extracted KMZ with %d other files\n", len(filesMap))
	for name := range filesMap {
		fmt.Printf("DEBUG: File in map: %s\n", name)
	}

	return kmlContent, filesMap, nil
}

// ExtractYoutubeURL extracts and converts a YouTube URL from a description
func ExtractYoutubeURL(description string) string {
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
			videoIDStart := startIdx + len(pattern)
			videoIDEnd := videoIDStart
			for videoIDEnd < len(description) && description[videoIDEnd] != ' ' &&
				description[videoIDEnd] != '"' &&
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

// ParseCoordinates parses KML coordinates format: "lng,lat[,altitude]"
func ParseCoordinates(coordStr string) (float64, float64, error) {
	coords := strings.Split(strings.TrimSpace(coordStr), ",")
	if len(coords) < 2 {
		return 0, 0, fmt.Errorf("invalid coordinates format")
	}

	lngRaw := strings.TrimSpace(coords[0])
	latRaw := strings.TrimSpace(coords[1])

	lng, err := strconv.ParseFloat(lngRaw, 64)
	if err != nil {
		return 0, 0, err
	}

	lat, err := strconv.ParseFloat(latRaw, 64)
	if err != nil {
		return 0, 0, err
	}

	// Return as (lat, lng) to match standard Go mapping expectations
	return lat, lng, nil

}
