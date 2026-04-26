-- migrations/004_create_places.up.sql

CREATE TABLE places (
  id              SERIAL PRIMARY KEY,
  city_id         INTEGER       NOT NULL 
                  REFERENCES cities(id) 
                  ON DELETE CASCADE,
  category_id     INTEGER       NOT NULL 
                  REFERENCES categories(id),
  name            TEXT  NOT NULL,
  subcategory     TEXT,
  -- e.g. 'ramen' under 'food'
  -- e.g. 'skyscraper' under 'architecture'
  address         TEXT,
  neighborhood    VARCHAR(100),
  -- critical for NYC — which borough/neighborhood
  borough         VARCHAR(50),
  -- NYC specific — Manhattan, Brooklyn, etc
  latitude        DECIMAL(9,6),
  longitude       DECIMAL(9,6),
  description     TEXT,
  -- editorial description
  price_tier      INTEGER,
  -- 0=free 1=cheap 2=moderate 3=expensive
  rating          DECIMAL(3,1),
  -- aggregated from sources
  hours           JSONB,
  -- flexible hours storage
  website         TEXT,
  phone           VARCHAR(50),

  -- Architecture specific fields
  -- NULL for non-architecture places
  architect       VARCHAR(200),
  year_built      INTEGER,
  architectural_style VARCHAR(100),
  -- links to style tag
  height_ft       INTEGER,
  floors          INTEGER,
  landmarked      BOOLEAN       DEFAULT false,
  
  -- Data provenance
  source          VARCHAR(50),
  -- overpass | foursquare | manual
  source_id       VARCHAR(200),
  -- original ID from source
  verified        BOOLEAN       DEFAULT false,
  -- manually verified by you
  
  created_at      TIMESTAMP     DEFAULT NOW(),
  updated_at      TIMESTAMP     DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_places_city_id 
  ON places(city_id);
CREATE INDEX idx_places_category_id 
  ON places(category_id);
CREATE INDEX idx_places_neighborhood 
  ON places(neighborhood);
CREATE INDEX idx_places_borough 
  ON places(borough);

-- Junction table — Places to Tags
-- One place can have many tags
CREATE TABLE place_tags (
  place_id    INTEGER NOT NULL 
              REFERENCES places(id) 
              ON DELETE CASCADE,
  tag_id      INTEGER NOT NULL 
              REFERENCES tags(id) 
              ON DELETE CASCADE,
  PRIMARY KEY (place_id, tag_id)
);

ALTER TABLE places 
ADD CONSTRAINT places_source_source_id_unique 
UNIQUE (source, source_id);
