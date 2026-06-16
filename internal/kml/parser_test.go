package kml

import (
	"testing"
)

func TestParseCoordinates(t *testing.T) {
	tests := []struct {
		name       string
		coordStr   string // KML format is "lng,lat"
		wantLat    float64
		wantLng    float64
		wantErr    bool
	}{
		{"valid_coords", "-0.09, 51.505", 51.505, -0.09, false},
		{"valid_coords_no_space", "-0.09,51.505", 51.505, -0.09, false},
		{"invalid_format", "51.505", 0, 0, true},
		{"missing_comma", "51.505-0.09", 0, 0, true},
		{"non_numeric", "abc,def", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lat, lng, err := ParseCoordinates(tt.coordStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCoordinates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if lat != tt.wantLat || lng != tt.wantLng {
					t.Errorf("ParseCoordinates() got (%v, %v), want (%v, %v)", lat, lng, tt.wantLat, tt.wantLng)
				}
			}
		})
	}
}

func TestExtractYoutubeURL(t *testing.T) {
	tests := []struct {
		name         string
		description  string
		want         string
	}{
		{
			"standard_watch_url",
			"Watch here: https://www.youtube.com/watch?v=dQw4w9WgXcQ123",
			"https://www.youtube.com/embed/dQw4w9WgXcQ123?modestbranding=1",
		},
		{
			"short_url",
			"Check this out youtu.be/dQw4w9WgXcQ",
			"https://www.youtube.com/embed/dQw4w9WgXcQ?modestbranding=1",
		},
		{
			"url_with_punctuation",
			"Video link: https://youtu.be/dQw4w9WgXcQ.",
			"https://www.youtube.com/embed/dQw4w9WgXcQ?modestbranding=1",
		},
		{
			"no_url",
			"No video here.",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractYoutubeURL(tt.description); got != tt.want {
				t.Errorf("ExtractYoutubeURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
