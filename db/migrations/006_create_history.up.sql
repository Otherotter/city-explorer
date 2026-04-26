-- migrations/006_create_history.up.sql

-- This is your history component
-- Append-only — rows are never updated
-- Every significant moment in a place 
-- or city gets recorded here
-- Tagged with pivot event tags

CREATE TABLE place_history (
  id            SERIAL PRIMARY KEY,
  place_id      INTEGER       NOT NULL 
                REFERENCES places(id) 
                ON DELETE CASCADE,
  city_id       INTEGER       NOT NULL 
                REFERENCES cities(id),
  event_date    DATE          NOT NULL,
  -- when this happened
  title         VARCHAR(200)  NOT NULL,
  -- short headline
  description   TEXT,
  -- what actually happened
  source_url    TEXT,
  -- link to article, reference, etc
  recorded_by   VARCHAR(50)   DEFAULT 'system',
  -- system or your name if manual
  created_at    TIMESTAMP     DEFAULT NOW()
  -- when YOU recorded this
);

-- Junction table — History entries to Tags
-- A pivot event can have multiple tags
CREATE TABLE history_tags (
  history_id  INTEGER NOT NULL 
              REFERENCES place_history(id) 
              ON DELETE CASCADE,
  tag_id      INTEGER NOT NULL 
              REFERENCES tags(id) 
              ON DELETE CASCADE,
  PRIMARY KEY (history_id, tag_id)
);

-- Index for time-based queries
-- Grafana will query this by date range
CREATE INDEX idx_place_history_event_date 
  ON place_history(event_date);
CREATE INDEX idx_place_history_place_id 
  ON place_history(place_id);
CREATE INDEX idx_place_history_city_id 
  ON place_history(city_id);

-- Seed some NYC history entries
INSERT INTO place_history (
  place_id, city_id, event_date, 
  title, description, 
  source_url, recorded_by
) VALUES
(
  (SELECT id FROM places 
   WHERE name = 'Chrysler Building'),
  1,
  '1930-05-27',
  'Chrysler Building Opens',
  'Completed in 1930 and briefly the 
   tallest building in the world at 
   1046 feet before the Empire State 
   Building surpassed it in 1931. 
   Designed by William Van Alen.',
  'https://en.wikipedia.org/wiki/Chrysler_Building',
  'manual'
),
(
  (SELECT id FROM places 
   WHERE name = 'High Line'),
  1,
  '2009-06-09',
  'High Line Section 1 Opens to Public',
  'First section of the elevated 
   rail park opened after years of 
   community advocacy. Transformed 
   the Meatpacking District and 
   West Chelsea permanently.',
  'https://www.thehighline.org/history',
  'manual'
),
(
  (SELECT id FROM places 
   WHERE name = 'Smorgasburg'),
  1,
  '2011-04-01',
  'Smorgasburg Founded',
  'Launched as a food-only offshoot 
   of the Brooklyn Flea market in 
   Williamsburg. Became the largest 
   open-air food market in America.',
  'https://www.smorgasburg.com/about',
  'manual'
);

-- Tag the history entries
INSERT INTO history_tags (history_id, tag_id)
VALUES
  (
    (SELECT id FROM place_history 
     WHERE title = 'Chrysler Building Opens'),
    (SELECT id FROM tags 
     WHERE slug = 'pivot-opening')
  ),
  (
    (SELECT id FROM place_history 
     WHERE title = 'High Line Section 1 Opens to Public'),
    (SELECT id FROM tags 
     WHERE slug = 'pivot-opening')
  ),
  (
    (SELECT id FROM place_history 
     WHERE title = 'High Line Section 1 Opens to Public'),
    (SELECT id FROM tags 
     WHERE slug = 'pivot-development')
  ),
  (
    (SELECT id FROM place_history 
     WHERE title = 'Smorgasburg Founded'),
    (SELECT id FROM tags 
     WHERE slug = 'pivot-opening')
  ),
  (
    (SELECT id FROM place_history 
     WHERE title = 'Smorgasburg Founded'),
    (SELECT id FROM tags 
     WHERE slug = 'pivot-cultural')
  );