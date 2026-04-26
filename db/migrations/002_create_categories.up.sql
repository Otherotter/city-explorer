-- migrations/002_create_categories.up.sql

-- Categories are their own table
-- This means adding a new category never 
-- requires a schema change
-- Just insert a new row

CREATE TABLE categories (
  id          SERIAL PRIMARY KEY,
  slug        VARCHAR(50)   NOT NULL UNIQUE,
  -- used in API routes
  label       VARCHAR(100)  NOT NULL,
  -- human readable
  description TEXT,
  icon        VARCHAR(50),
  -- emoji or icon name for frontend
  created_at  TIMESTAMP     DEFAULT NOW()
);

-- Seed all categories including Architecture
INSERT INTO categories (slug, label, description, icon) 
VALUES
  (
    'food',
    'Food & Drink',
    'Local restaurants, cafes, bars, 
     markets, and hidden gems',
    '🍜'
  ),
  (
    'attraction',
    'Must-See Attractions',
    'Landmarks, museums, neighborhoods, 
     and experiences locals recommend',
    '🏛️'
  ),
  (
    'architecture',
    'Architecture',
    'Notable buildings, structural landmarks,
     design districts, and urban form',
    '🏗️'
  ),
  (
    'thrift',
    'Thrift & Vintage',
    'Second-hand shops, vintage stores, 
     record shops, and bookstores',
    '🛍️'
  ),
  (
    'events',
    'Events',
    'Things happening this week — 
     music, food, art, culture',
    '📅'
  ),
  (
    'nature',
    'Nature & Outdoors',
    'Parks, trails, waterfronts, 
     gardens, and viewpoints',
    '🌿'
  ),
  (
    'study',
    'Study Spots',
    'Coffee shops, libraries, and 
     coworking spaces with good wifi',
    '☕'
  ),
  (
    'social',
    'Meet People',
    'Meetup groups, community events, 
     communal dining, and social spaces',
    '🤝'
  ),
  (
    'climbing',
    'Climbing & Bouldering',
    'Indoor gyms, outdoor crags, 
     and bouldering spots',
    '🧗'
  );