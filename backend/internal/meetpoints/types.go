package meetpoints

import "time"

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type MeetRequest struct {
	OriginA string `json:"originA"`
	OriginB string `json:"originB"`
}

type MeetResponse struct {
	Origins      Origins     `json:"origins"`
	Best         MeetPoint   `json:"best"`
	Alternatives []MeetPoint `json:"alternatives,omitempty"`
	Debug        Debug       `json:"debug,omitempty"`
}

type Origins struct {
	A Origin `json:"a"`
	B Origin `json:"b"`
}

type Origin struct {
	Address string `json:"address"`
	Point   LatLng `json:"point"`
}

type MeetPoint struct {
	Point           LatLng `json:"point"`
	ETAFromA        int    `json:"etaASeconds"`
	ETAFromB        int    `json:"etaBSeconds"`
	MaxETASeconds   int    `json:"maxEtaSeconds"`
	DiffSeconds     int    `json:"diffSeconds"`
	DistanceAMeters int    `json:"distanceAMeters"`
	DistanceBMeters int    `json:"distanceBMeters"`
}

type Debug struct {
	Midpoint LatLng `json:"midpoint"`
}

type Geocoder interface {
	Geocode(address string) (LatLng, error)
}

type Router interface {
	RouteETA(origin, dest LatLng, departAt time.Time) (durationSeconds int, distanceMeters int, err error)
}
