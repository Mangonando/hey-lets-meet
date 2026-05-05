package meetpoints

import (
	"hash/fnv"
	"math"
	"strconv"
	"time"
)

type MockRouter struct {
	WalkingSpeedMps float64
}

func (m MockRouter) RouteETA(origin, dest LatLng, _ time.Time) (int, int, error) {
	walkingSpeed := m.WalkingSpeedMps
	if walkingSpeed <= 0 {
		walkingSpeed = 1.4
	}

	distanceMeters := HaversineMeters(origin, dest)
	durationSeconds := distanceMeters / walkingSpeed

	jitterSeconds := jitter(origin, dest)
	durationSeconds += float64(jitterSeconds)

	return int(math.Round(durationSeconds)), int(math.Round(distanceMeters)), nil
}

func jitter(origin, dest LatLng) int {
	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(
		formatFloat(origin.Lat) + "," + formatFloat(origin.Lng) + "|" + formatFloat(dest.Lat) + "," + formatFloat(dest.Lng),
	))
	return int(hasher.Sum32() % 30)
}

func formatFloat(val float64) string {
	return strconv.FormatFloat(val, 'f', 6, 64)
}
