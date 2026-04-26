-- migrations/001_create_cities.up.sql

CREATE TABLE cities (
  id              SERIAL PRIMARY KEY,
  name            VARCHAR(100)    NOT NULL,
  country         VARCHAR(100)    NOT NULL,
  region          VARCHAR(100),
  -- state, province, county
  latitude        DECIMAL(9,6)    NOT NULL,
  longitude       DECIMAL(9,6)    NOT NULL,
  population      INTEGER,
  timezone        VARCHAR(50),
  -- e.g. America/New_York
  description     TEXT,
  -- short editorial note about the city
  created_at      TIMESTAMP       DEFAULT NOW(),
  updated_at      TIMESTAMP       DEFAULT NOW()
);

-- Index for name lookups
CREATE INDEX idx_cities_name 
  ON cities(name);

-- Seed NYC immediately
INSERT INTO cities (
  name, 
  country, 
  region, 
  latitude, 
  longitude, 
  population, 
  timezone, 
  description
) VALUES (
  'New York City',
  'United States',
  'New York',
  40.712776,
  -74.005974,
  8336817,
  'America/New_York',
  'The city that never sleeps. 
   Five boroughs, infinite neighborhoods, 
   every culture on earth.'
);