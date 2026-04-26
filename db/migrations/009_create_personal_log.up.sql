-- migrations/009_create_personal_log.up.sql

CREATE TABLE personal_log (
  id              SERIAL PRIMARY KEY,
  place_id        INTEGER 
                  REFERENCES places(id),
  -- nullable — can log a city 
  -- without a specific place
  city_id         INTEGER       NOT NULL 
                  REFERENCES cities(id),
  visited_at      TIMESTAMP     NOT NULL,
  personal_rating INTEGER 
                  CHECK (
                    personal_rating >= 1 
                    AND personal_rating <= 10
                  ),
  notes           TEXT,
  would_return    BOOLEAN,
  mood            VARCHAR(50),
  -- how you felt that day
  -- good context for patterns
  who_with        VARCHAR(100),
  -- solo | friends | family | date
  logged_at       TIMESTAMP     DEFAULT NOW()
);

-- Junction table — Personal Log to Tags
CREATE TABLE personal_log_tags (
  log_id    INTEGER NOT NULL 
            REFERENCES personal_log(id) 
            ON DELETE CASCADE,
  tag_id    INTEGER NOT NULL 
            REFERENCES tags(id) 
            ON DELETE CASCADE,
  PRIMARY KEY (log_id, tag_id)
);

CREATE INDEX idx_personal_log_city_id 
  ON personal_log(city_id);
CREATE INDEX idx_personal_log_visited_at 
  ON personal_log(visited_at);