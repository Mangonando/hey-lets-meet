package meetpoints

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ORSConfig struct {
	BaseURL string
	APIKey  string
	Timeout time.Duration
}

type ORSGeocoder struct {
	HTTP   *http.Client
	Config ORSConfig
	Cache  CacheRepo
}

type ORSRouter struct {
	HTTP   *http.Client
	Config ORSConfig
	Cache  CacheRepo
}

func normalizeAddress(a string) string {
	return strings.TrimSpace(strings.ToLower(a))
}

func (geocoder ORSGeocoder) Geocode(address string) (LatLng, error) {
	norm := normalizeAddress(address)
	if norm == "" {
		return LatLng{}, errors.New("empty address")
	}

	if point, ok, err := geocoder.Cache.GetGeocode(norm); err != nil {
		return LatLng{}, err
	} else if ok {
		return point, nil
	}

	parsedURL, _ := url.Parse(strings.TrimRight(geocoder.Config.BaseURL, "/") + "/geocode/search")
	queryParams := parsedURL.Query()
	queryParams.Set("api_key", geocoder.Config.APIKey)
	queryParams.Set("text", address)
	queryParams.Set("size", "1")
	parsedURL.RawQuery = queryParams.Encode()

	req, err := http.NewRequest(http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return LatLng{}, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := geocoder.HTTP.Do(req)
	if err != nil {
		return LatLng{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return LatLng{}, fmt.Errorf("geocode failed: HTTP %d", resp.StatusCode)
	}

	var geocodeResponse struct {
		Features []struct {
			Geometry struct {
				Coordinates []float64 `json:"coordinates"`
			} `json:"geometry"`
		} `json:"features"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&geocodeResponse); err != nil {
		return LatLng{}, err
	}
	if len(geocodeResponse.Features) == 0 || len(geocodeResponse.Features[0].Geometry.Coordinates) < 2 {
		return LatLng{}, errors.New("no geocoding result")
	}

	lng := geocodeResponse.Features[0].Geometry.Coordinates[0]
	lat := geocodeResponse.Features[0].Geometry.Coordinates[1]
	point := LatLng{Lat: lat, Lng: lng}

	_ = geocoder.Cache.PutGeocode(norm, point)
	return point, nil
}

func (router ORSRouter) RouteETA(origin, dest LatLng, _ time.Time) (durationSeconds int, distanceMeters int, err error) {
	key := fmt.Sprintf(
		"walk:%s:%s",
		roundCoordKey(origin),
		roundCoordKey(dest),
	)

	if durationSeconds, distanceMeters, ok, err := router.Cache.GetWalkRoute(key); err != nil {
		return 0, 0, err
	} else if ok {
		return durationSeconds, distanceMeters, nil
	}

	endpoint := strings.TrimRight(router.Config.BaseURL, "/") + "/v2/directions/foot-walking/json"

	requestBody := map[string]any{
		"coordinates": [][]float64{
			{origin.Lng, origin.Lat},
			{dest.Lng, dest.Lat},
		},
		"instructions": false,
	}
	b, _ := json.Marshal(requestBody)

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(b))
	if err != nil {
		return 0, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", router.Config.APIKey)

	resp, err := router.HTTP.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, 0, fmt.Errorf("directions failed: HTTP %d", resp.StatusCode)
	}

	var routeResponse struct {
		Routes []struct {
			Summary struct {
				Distance float64 `json:"distance"`
				Duration float64 `json:"duration"`
			} `json:"summary"`
		} `json:"routes"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&routeResponse); err != nil {
		return 0, 0, err
	}
	if len(routeResponse.Routes) == 0 {
		return 0, 0, errors.New("no route returned")
	}

	distanceMeters = int(routeResponse.Routes[0].Summary.Distance + 0.5)
	durationSeconds = int(routeResponse.Routes[0].Summary.Duration + 0.5)

	_ = router.Cache.PutWalkRoute(key, durationSeconds, distanceMeters)
	return durationSeconds, distanceMeters, nil
}

func roundCoordKey(p LatLng) string {
	return fmt.Sprintf("%.5f,%.5f", p.Lat, p.Lng)
}
