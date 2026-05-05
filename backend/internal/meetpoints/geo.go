package meetpoints

import "math"

// average radius of the earth in meters
const earthRadiusMeters = 6371000.0

// straight line distance two coordinates on earth's surface
func HaversineMeters(a, b LatLng) float64 {
	lat1 := deg2rad(a.Lat)
	lng1 := deg2rad(a.Lng)
	lat2 := deg2rad(b.Lat)
	lng2 := deg2rad(b.Lng)

	dlat := lat2 - lat1
	dlng := lng2 - lng1

	sinDLat := math.Sin(dlat / 2)
	sinDLng := math.Sin(dlng / 2)

	// result of the haversine formula. value between 2 points in a circle
	haversine := sinDLat*sinDLat + math.Cos(lat1)*math.Cos(lat2)*sinDLng*sinDLng
	// angle at earth's center between 2 points
	centralAngle := 2 * math.Atan2(math.Sqrt(haversine), math.Sqrt(1-haversine))

	return earthRadiusMeters * centralAngle
}

// geographic center between 2 coordinates
func Midpoint(a, b LatLng) LatLng {
	return LatLng{
		Lat: (a.Lat + b.Lat) / 2,
		Lng: (a.Lng + b.Lng) / 2,
	}
}

// shift coordinate b x meters east and y meters north
func OffsetMeters(center LatLng, eastMeters, northMeters float64) LatLng {
	// one degree of latitude = 111,320
	const metersPerDegLat = 111320.0
	latDelta := northMeters / metersPerDegLat
	lngDelta := eastMeters / (metersPerDegLat * math.Cos(deg2rad(center.Lat)))

	return LatLng{
		Lat: center.Lat + latDelta,
		Lng: center.Lng + lngDelta,
	}
}

// math.Pi converst degrees to radians
func deg2rad(degrees float64) float64 { return degrees * math.Pi / 180.0 }
