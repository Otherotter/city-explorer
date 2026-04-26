-- migrations/007_create_events.up.sql

CREATE TABLE events (
  id            SERIAL PRIMARY KEY,
  city_id       INTEGER       NOT NULL 
                REFERENCES cities(id) 
                ON DELETE CASCADE,
  name          VARCHAR(200)  NOT NULL,
  category_id   INTEGER       
                REFERENCES categories(id),
  description   TEXT,
  venue         VARCHAR(200),
  neighborhood  VARCHAR(100),
  borough       VARCHAR(50),
  latitude      DECIMAL(9,6),
  longitude     DECIMAL(9,6),
  start_time    TIMESTAMP     NOT NULL,
  end_time      TIMESTAMP,
  is_free       BOOLEAN       DEFAULT false,
  price_min     DECIMAL(8,2),
  url           TEXT,
  source        VARCHAR(50),
  source_id     VARCHAR(200),
  created_at    TIMESTAMP     DEFAULT NOW(),
  updated_at    TIMESTAMP     DEFAULT NOW()
);

CREATE INDEX idx_events_city_id 
  ON events(city_id);
CREATE INDEX idx_events_start_time 
  ON events(start_time);
CREATE INDEX idx_events_is_free 
  ON events(is_free);