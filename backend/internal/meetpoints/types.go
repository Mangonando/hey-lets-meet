package meetpoints

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type SuggestRequest struct {
	OriginA string `json:"originA"`
	OriginB string `json:"originB"`
}
