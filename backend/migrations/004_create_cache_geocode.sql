CREATE TABLE IF NOT EXISTS cache_geocode (
  address_norm TEXT PRIMARY KEY,
  lat REAL NOT NULL,
  lng REAL NOT NULL,
  cached_at TEXT NOT NULL DEFAULT (datetime('now'))
);