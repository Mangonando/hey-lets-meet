package meetpoints

import (
	"testing"
	"time"
)

func TestSuggestDeterministic(t *testing.T) {
	svc := &Service{
		Geocoder: MockGeocoder{},
		Router:   MockRouter{WalkingSpeedMps: 1.4},
	}

	resp, err := svc.Suggest(MeetRequest{
		OriginA: "Alexanderplatz",
		OriginB: "Hermannplatz",
	}, time.Unix(0, 0), DefaultOptions())
	if err != nil {
		t.Fatalf("Suggest error: %v", err)
	}

	if resp.Best.MaxETASeconds <= 0 {
		t.Fatalf("expected positive ETA, got %d", resp.Best.MaxETASeconds)
	}
	if len(resp.Alternatives) == 0 {
		t.Fatalf("expected alternatives")
	}
}
