package meetpoints

import (
	"math"
	"sort"
	"time"
)

type Service struct {
	Geocoder Geocoder
	Router   Router
}

type Options struct {
	GridSize     int
	RadiusMeters float64
}

func DefaultOptions() Options {
	return Options{GridSize: 3, RadiusMeters: 2000}
}

func (s *Service) Suggest(req MeetRequest, departAt time.Time, options Options) (MeetResponse, error) {
	pointA, err := s.Geocoder.Geocode(req.OriginA)
	if err != nil {
		return MeetResponse{}, err
	}
	pointB, err := s.Geocoder.Geocode(req.OriginB)
	if err != nil {
		return MeetResponse{}, err
	}

	midpoint := Midpoint(pointA, pointB)
	candidates := gridCandidates(midpoint, options.GridSize, options.RadiusMeters)

	results := make([]MeetPoint, 0, len(candidates))
	for _, candidate := range candidates {
		etaASec, distAMeters, err := s.Router.RouteETA(pointA, candidate, departAt)
		if err != nil {
			return MeetResponse{}, err
		}
		etaBSec, distBMeters, err := s.Router.RouteETA(pointB, candidate, departAt)
		if err != nil {
			return MeetResponse{}, err
		}

		maxETA := maxInt(etaASec, etaBSec)
		diff := int(math.Abs(float64(etaASec - etaBSec)))

		results = append(results, MeetPoint{
			Point:           candidate,
			ETAFromA:        etaASec,
			ETAFromB:        etaBSec,
			MaxETASeconds:   maxETA,
			DiffSeconds:     diff,
			DistanceAMeters: distAMeters,
			DistanceBMeters: distBMeters,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].MaxETASeconds != results[j].MaxETASeconds {
			return results[i].MaxETASeconds < results[j].MaxETASeconds
		}
		if results[i].DiffSeconds != results[j].DiffSeconds {
			return results[i].DiffSeconds < results[j].DiffSeconds
		}
		totalDistI := results[i].DistanceAMeters + results[i].DistanceBMeters
		totalDistJ := results[j].DistanceAMeters + results[j].DistanceBMeters
		return totalDistI < totalDistJ
	})

	best := results[0]
	alternatives := []MeetPoint{}
	if len(results) > 1 {
		// return top 3 alternatives excluding the best one
		limit := 4
		if len(results) < limit {
			limit = len(results)
		}
		alternatives = append(alternatives, results[1:limit]...)
	}

	return MeetResponse{
		Origins: Origins{
			A: Origin{Address: req.OriginA, Point: pointA},
			B: Origin{Address: req.OriginB, Point: pointB},
		},
		Best:         best,
		Alternatives: alternatives,
		Debug:        Debug{Midpoint: midpoint},
	}, nil
}

func gridCandidates(center LatLng, gridSize int, radiusMeters float64) []LatLng {
	if gridSize < 3 {
		gridSize = 3
	}
	if gridSize%2 == 0 {
		gridSize++
	}
	steps := gridSize - 1
	step := (2 * radiusMeters) / float64(steps)

	start := -radiusMeters
	out := make([]LatLng, 0, gridSize*gridSize)

	for row := 0; row < gridSize; row++ {
		north := start + float64(row)*step
		for col := 0; col < gridSize; col++ {
			east := start + float64(col)*step
			out = append(out, OffsetMeters(center, east, north))
		}
	}
	return out
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
