package meetpoints

import "database/sql"

type CacheRepo struct {
	DB *sql.DB
}

func (repo CacheRepo) GetGeocode(addressNorm string) (LatLng, bool, error) {
	row := repo.DB.QueryRow(`SELECT lat, lng FROM cache_geocode WHERE address_norm = ?`, addressNorm)
	var lat, lng float64
	if err := row.Scan(&lat, &lng); err != nil {
		if err == sql.ErrNoRows {
			return LatLng{}, false, nil
		}
		return LatLng{}, false, err
	}
	return LatLng{Lat: lat, Lng: lng}, true, nil
}

func (repo CacheRepo) PutGeocode(addressNorm string, point LatLng) error {
	_, err := repo.DB.Exec(`
INSERT INTO cache_geocode(address_norm, lat, lng, cached_at)
VALUES (?, ?, ?, datetime('now'))
ON CONFLICT(address_norm) DO UPDATE SET
  lat=excluded.lat,
  lng=excluded.lng,
  cached_at=datetime('now')
`, addressNorm, point.Lat, point.Lng)
	return err
}

func (repo CacheRepo) GetWalkRoute(key string) (durationSec int, distanceMeters int, ok bool, err error) {
	row := repo.DB.QueryRow(`SELECT duration_seconds, distance_meters FROM cache_route_walk WHERE key = ?`, key)
	if err := row.Scan(&durationSec, &distanceMeters); err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, false, nil
		}
		return 0, 0, false, err
	}
	return durationSec, distanceMeters, true, nil
}

func (repo CacheRepo) PutWalkRoute(key string, durationSec int, distanceMeters int) error {
	_, err := repo.DB.Exec(`
INSERT INTO cache_route_walk(key, duration_seconds, distance_meters, cached_at)
VALUES (?, ?, ?, datetime('now'))
ON CONFLICT(key) DO UPDATE SET
  duration_seconds=excluded.duration_seconds,
  distance_meters=excluded.distance_meters,
  cached_at=datetime('now')
`, key, durationSec, distanceMeters)
	return err
}
