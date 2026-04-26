-- Seed NYC Places across all categories

INSERT INTO places (
  city_id, category_id, name, subcategory,
  address, neighborhood, borough,
  latitude, longitude, description,
  price_tier, rating, source, verified
) VALUES

-- ─────────────────────────────────────────
-- FOOD
-- ─────────────────────────────────────────
(
  1, -- NYC
  (SELECT id FROM categories WHERE slug = 'food'),
  'Xi''an Famous Foods',
  'noodles',
  '81 St Marks Pl, New York, NY 10003',
  'East Village',
  'Manhattan',
  40.728000, -73.983000,
  'Hand-ripped noodles and cumin lamb 
   that will change how you think about 
   Chinese food. Order the spicy 
   tingly lamb face salad.',
  1, 4.5, 'manual', true
),
(
  1,
  (SELECT id FROM categories WHERE slug = 'food'),
  'Smorgasburg',
  'food market',
  'East River State Park, Brooklyn',
  'Williamsburg',
  'Brooklyn',
  40.715000, -73.963000,
  'Open-air food market every Saturday. 
   100 local vendors. The best eating 
   day you can have in NYC for under $30.',
  1, 4.7, 'manual', true
),
(
  1,
  (SELECT id FROM categories WHERE slug = 'food'),
  'Superiority Burger',
  'vegetarian',
  '119 Avenue A, New York, NY 10009',
  'East Village',
  'Manhattan',
  40.726000, -73.981000,
  'Do not let the vegetarian label 
   fool you. The best burger in NYC 
   is meat-free. Cash only. 
   Lines around the block.',
  1, 4.6, 'manual', true
),

-- ─────────────────────────────────────────
-- ARCHITECTURE
-- ─────────────────────────────────────────
(
  1,
  (SELECT id FROM categories WHERE slug = 'architecture'),
  'Chrysler Building',
  'skyscraper',
  '405 Lexington Ave, New York, NY 10174',
  'Midtown',
  'Manhattan',
  40.751717, -73.975504,
  'The finest Art Deco skyscraper ever 
   built. The eagle gargoyles at the 
   31st floor are the best detail in 
   NYC architecture. Look up.',
  0, 4.9, 'manual', true
),
(
  1,
  (SELECT id FROM categories WHERE slug = 'architecture'),
  'Flatiron Building',
  'skyscraper',
  '175 Fifth Ave, New York, NY 10010',
  'Flatiron',
  'Manhattan',
  40.741061, -73.989699,
  'The original NYC icon. Built in 1902 
   on a triangular plot. Walk the full 
   perimeter — it looks different 
   from every angle.',
  0, 4.8, 'manual', true
),
(
  1,
  (SELECT id FROM categories WHERE slug = 'architecture'),
  'Woolworth Building',
  'skyscraper',
  '233 Broadway, New York, NY 10279',
  'Tribeca',
  'Manhattan',
  40.712775, -74.008057,
  'Gothic Revival skyscraper from 1913. 
   Called the Cathedral of Commerce 
   when it was built. The lobby 
   is open to the public.',
  0, 4.7, 'manual', true
),
(
  1,
  (SELECT id FROM categories WHERE slug = 'architecture'),
  'High Line',
  'urban design',
  'Gansevoort St to 34th St, Manhattan',
  'Chelsea',
  'Manhattan',
  40.748000, -74.004000,
  'Elevated railway converted to linear 
   park. Study how the design team 
   kept the original rail tracks 
   visible throughout.',
  0, 4.6, 'manual', true
),

-- ─────────────────────────────────────────
-- NATURE
-- ─────────────────────────────────────────
(
  1,
  (SELECT id FROM categories WHERE slug = 'nature'),
  'Inwood Hill Park',
  'park',
  'Inwood Hill Park, New York, NY 10034',
  'Inwood',
  'Manhattan',
  40.873000, -73.921000,
  'The last remaining wild forest 
   in Manhattan. Ancient caves, 
   herons, and almost no tourists. 
   The Manhattan most people 
   never see.',
  0, 4.8, 'manual', true
),
(
  1,
  (SELECT id FROM categories WHERE slug = 'nature'),
  'Jamaica Bay Wildlife Refuge',
  'wildlife',
  'Cross Bay Blvd, Broad Channel, NY',
  'Broad Channel',
  'Queens',
  40.619000, -73.830000,
  'Bird sanctuary inside NYC limits. 
   330+ species spotted. Take the A 
   train to Broad Channel. 
   Completely free.',
  0, 4.7, 'manual', true
),

-- ─────────────────────────────────────────
-- THRIFT
-- ─────────────────────────────────────────
(
  1,
  (SELECT id FROM categories WHERE slug = 'thrift'),
  'Housing Works Thrift Shop',
  'thrift store',
  '143 W 17th St, New York, NY 10011',
  'Chelsea',
  'Manhattan',
  40.740000, -73.997000,
  'Best thrift chain in NYC. 
   Proceeds go to HIV/AIDS services. 
   The Chelsea location gets 
   incredible donations.',
  1, 4.5, 'manual', true
),
(
  1,
  (SELECT id FROM categories WHERE slug = 'thrift'),
  'Academy Records',
  'record store',
  '85 Oak St, Brooklyn, NY 11222',
  'Greenpoint',
  'Brooklyn',
  40.724000, -73.951000,
  'Used and new vinyl. Jazz section 
   is exceptional. The staff 
   actually knows music.',
  1, 4.7, 'manual', true
),

-- ─────────────────────────────────────────
-- STUDY SPOTS
-- ─────────────────────────────────────────
(
  1,
  (SELECT id FROM categories WHERE slug = 'study'),
  'New York Public Library — Rose Main Reading Room',
  'library',
  '476 5th Ave, New York, NY 10018',
  'Midtown',
  'Manhattan',
  40.753200, -73.982300,
  'The most beautiful room to work 
   in NYC. 78-foot ceilings, 
   original chandeliers, absolute 
   silence. Free. No reservation needed.',
  0, 4.9, 'manual', true
),
(
  1,
  (SELECT id FROM categories WHERE slug = 'study'),
  'Cafe Grumpy',
  'cafe',
  '224 W 20th St, New York, NY 10011',
  'Chelsea',
  'Manhattan',
  40.742000, -74.000000,
  'Serious coffee, good wifi, 
   no laptop shaming. 
   Quieter than most Chelsea cafes.',
  1, 4.4, 'manual', true
),

-- ─────────────────────────────────────────
-- ATTRACTIONS
-- ─────────────────────────────────────────
(
  1,
  (SELECT id FROM categories WHERE slug = 'attraction'),
  'The Tenement Museum',
  'museum',
  '103 Orchard St, New York, NY 10002',
  'Lower East Side',
  'Manhattan',
  40.718000, -73.990000,
  'The most honest museum in NYC. 
   Preserved immigrant apartments 
   from the 1860s to 1930s. 
   Book a tour — worth every penny.',
  2, 4.9, 'manual', true
),
(
  1,
  (SELECT id FROM categories WHERE slug = 'attraction'),
  'Flushing Meadows Corona Park',
  'park and landmark',
  'Flushing Meadows, Queens, NY',
  'Flushing',
  'Queens',
  40.730000, -73.842000,
  'Home of the 1964 World''s Fair. 
   The Unisphere is still there. 
   Queens Museum is underrated. 
   Completely off the tourist trail.',
  0, 4.5, 'manual', true
),

-- ─────────────────────────────────────────
-- SOCIAL
-- ─────────────────────────────────────────
(
  1,
  (SELECT id FROM categories WHERE slug = 'social'),
  'Hester Street Fair',
  'market and social',
  'Hester St and Essex St, Manhattan',
  'Lower East Side',
  'Manhattan',
  40.715000, -73.990000,
  'Weekend market with makers, 
   food vendors, and a genuinely 
   local crowd. Good place to 
   talk to strangers without 
   it being weird.',
  0, 4.4, 'manual', true
);