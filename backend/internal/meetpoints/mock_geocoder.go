package meetpoints

import (
	"errors"
	"hash/fnv"
	"math"
	"strings"
)

type MockGeocoder struct{}

func (m MockGeocoder) Geocode(address string) (LatLng, error) {
	normalized := strings.TrimSpace(strings.ToLower(address))
	if normalized == "" {
		return LatLng{}, errors.New("empty address")
	}

	switch normalized {
	case "alexanderplatz", "alexanderplatz berlin":
		return LatLng{Lat: 52.521918, Lng: 13.413215}, nil
	case "hermannplatz", "hermannplatz berlin":
		return LatLng{Lat: 52.486355, Lng: 13.424318}, nil
	case "brandenburger tor", "brandenburg gate", "brandenburger tor berlin":
		return LatLng{Lat: 52.516275, Lng: 13.377704}, nil
	}

	// Stable fallback: hash address to a point around central Berlin.
	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(normalized))
	hashValue := float64(hasher.Sum32())

	// Map to small offsets: more or less 6km
	latOffset := (fract(hashValue/1000.0) - 0.5) * 0.10
	lngOffset := (fract(hashValue/7000.0) - 0.5) * 0.16

	berlinCenter := LatLng{Lat: 52.5200, Lng: 13.4050}
	return LatLng{
		Lat: berlinCenter.Lat + latOffset,
		Lng: berlinCenter.Lng + lngOffset,
	}, nil
}

func fract(x float64) float64 {
	return x - math.Floor(x)
}
