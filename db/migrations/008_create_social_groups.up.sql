-- migrations/008_create_social_groups.up.sql

CREATE TABLE social_groups (
  id            SERIAL PRIMARY KEY,
  city_id       INTEGER       NOT NULL 
                REFERENCES cities(id) 
                ON DELETE CASCADE,
  name          VARCHAR(200)  NOT NULL,
  category      VARCHAR(100),
  -- climbing | music | tech | food
  description   TEXT,
  meeting_cadence VARCHAR(100),
  -- weekly | monthly | irregular
  url           TEXT,
  member_count  INTEGER,
  neighborhood  VARCHAR(100),
  borough       VARCHAR(50),
  source        VARCHAR(50),
  created_at    TIMESTAMP     DEFAULT NOW(),
  updated_at    TIMESTAMP     DEFAULT NOW()
);

-- Seed NYC climbing community
INSERT INTO social_groups (
  city_id, name, category, 
  description, meeting_cadence, 
  url, source
) VALUES (
  1,
  'NYC Bouldering Collective',
  'climbing',
  'Community of boulderers across 
   NYC gyms — The Cliffs, Brooklyn 
   Boulders, Vital. Regular meetups 
   and outdoor trips to the 
   Gunks and Bishop.',
  'weekly',
  'https://www.meetup.com',
  'manual'
);