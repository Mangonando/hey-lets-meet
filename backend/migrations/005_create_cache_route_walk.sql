CREATE TABLE IF NOT EXISTS cache_route_walk (
  key TEXT PRIMARY KEY,
  duration_seconds INTEGER NOT NULL,
  distance_meters INTEGER NOT NULL,
  cached_at TEXT NOT NULL DEFAULT (datetime('now'))
);