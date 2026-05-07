package meetpoints

import "testing"

func TestGridCandidatesCount(t *testing.T) {
	center := LatLng{Lat: 52.52, Lng: 13.405}

	candidates := gridCandidates(center, 7, 2000)
	if got, want := len(candidates), 49; got != want {
		t.Fatalf("candidates = %d, want %d", got, want)
	}

	foundNearCenter := false
	for _, candidate := range candidates {
		if HaversineMeters(center, candidate) < 2.0 {
			foundNearCenter = true
			break
		}
	}
	if !foundNearCenter {
		t.Fatalf("expected a candidate near center")
	}
}
