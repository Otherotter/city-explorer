-- migrations/003_create_tags.up.sql

-- Tags are flexible labels that can be 
-- applied to anything
-- Places, events, history entries, personal logs
-- This is where your pivot event marking lives

CREATE TABLE tags (
  id          SERIAL PRIMARY KEY,
  slug        VARCHAR(100)  NOT NULL UNIQUE,
  label       VARCHAR(100)  NOT NULL,
  color       VARCHAR(7),
  -- hex color for frontend display
  created_at  TIMESTAMP     DEFAULT NOW()
);

-- Seed initial tags
-- General place tags
INSERT INTO tags (slug, label, color) VALUES
  ('cash-only',         'Cash Only',          '#F59E0B'),
  ('hidden-gem',        'Hidden Gem',          '#10B981'),
  ('tourist-trap',      'Tourist Trap',        '#EF4444'),
  ('local-favorite',    'Local Favorite',      '#3B82F6'),
  ('open-late',         'Open Late',           '#8B5CF6'),
  ('outdoor-seating',   'Outdoor Seating',     '#06B6D4'),
  ('free-entry',        'Free Entry',          '#10B981'),
  ('wifi-good',         'Good Wifi',           '#3B82F6'),
  ('quiet',             'Quiet',               '#6B7280'),
  ('instagrammable',    'Instagrammable',      '#EC4899'),
  ('dog-friendly',      'Dog Friendly',        '#F59E0B'),

-- Architecture specific tags
  ('art-deco',          'Art Deco',            '#D97706'),
  ('brutalist',         'Brutalist',           '#4B5563'),
  ('beaux-arts',        'Beaux Arts',          '#7C3AED'),
  ('modernist',         'Modernist',           '#0EA5E9'),
  ('gothic-revival',    'Gothic Revival',      '#1D4ED8'),
  ('cast-iron',         'Cast Iron',           '#92400E'),
  ('landmarked',        'NYC Landmarked',      '#DC2626'),

-- History / Pivot Event tags
-- These mark significant moments in 
-- a place or city over time
  ('pivot-opening',     'Grand Opening',       '#10B981'),
  ('pivot-closing',     'Permanently Closed',  '#EF4444'),
  ('pivot-renovated',   'Recently Renovated',  '#F59E0B'),
  ('pivot-ownership',   'New Ownership',       '#8B5CF6'),
  ('pivot-disaster',    'Major Event',         '#DC2626'),
  ('pivot-cultural',    'Cultural Shift',      '#EC4899'),
  ('pivot-development', 'Urban Development',   '#0EA5E9'),
  ('pivot-recognition', 'Award or Recognition','#D97706'),

-- Personal tags
  ('been-here',         'Been Here',           '#10B981'),
  ('want-to-go',        'Want to Go',          '#3B82F6'),
  ('would-return',      'Would Return',        '#8B5CF6'),
  ('overrated',         'Overrated',           '#EF4444');