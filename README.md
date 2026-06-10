# Train Cam Tracker

A Wails-based desktop application for tracking train cameras and viewing live feeds from multiple locations.

## Features

- **Interactive Map**: View all train cameras on a Leaflet map
- **Live YouTube Feeds**: Watch live camera feeds directly in the application
- **KML/KMZ Import**: Import camera locations from Google Maps KML/KMZ files

## KML/KMZ Import

The application now supports importing camera locations from KML (Keyhole Markup Language) and KMZ (compressed KML) files exported from Google Maps.

### How to Use

1. **Create or Export a KML File**:
   - In Google Maps, create place markers at your camera locations
   - Export the map as a KML file using Google My Maps
   - Or use any tool that exports KML/KMZ format

2. **Import into Train Cam Tracker**:
   - Click the "Import KML/KMZ" button in the sidebar
   - Select your `.kml` or `.kmz` file
   - The application will extract and display the camera locations

3. **Camera Information**:
   - **Name**: Extracted from the placemark name
   - **Location**: Extracted from the placemark coordinates (latitude/longitude)
   - **Description**: Extracted from the placemark description field
   - **YouTube URL**: Can be manually added after import

### KML File Format

The application supports standard KML files with the following structure:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
  <Document>
    <Placemark>
      <name>Camera Name</name>
      <description>Optional description</description>
      <Point>
        <coordinates>longitude,latitude,altitude</coordinates>
      </Point>
    </Placemark>
  </Document>
</kml>
```

**Note**: KML coordinates are in the format `longitude,latitude` (not `latitude,longitude`), with an optional altitude component.

### Example

A sample KML file is provided in `sample.kml` showing camera locations around London.

## Building

```bash
wails build
```

## Development

```bash
wails dev
```

## Architecture

### Backend (Go)

- `app.go`: Main application logic with camera management and KML import functionality
- `main.go`: Wails application setup

### Frontend (JavaScript)

- `frontend/index.html`: Application UI with map container and sidebar
- `frontend/main.js`: Map initialization, marker management, and import handling
- `frontend/style.css`: Styling for the application

## Dependencies

- **Go**: Backend runtime
- **Wails v2**: Desktop framework
- **Leaflet.js**: Interactive map library
- **OpenStreetMap**: Map tile provider
- **YouTube**: Live feed embeds

## Adding YouTube URLs

After importing cameras from a KML file, you can add YouTube URLs to display live feeds:

1. Select a camera by clicking its map marker
2. The sidebar will show the camera details
3. To add a YouTube feed, edit the `YoutubeURL` field in the application

## Limitations

- Each camera location requires exactly one Point coordinate
- YouTube URLs must be in embed format: `https://www.youtube.com/embed/{VIDEO_ID}`
- Placemarks without Point coordinates will be skipped during import
