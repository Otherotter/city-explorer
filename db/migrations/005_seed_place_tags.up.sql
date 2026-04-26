-- migrations/005_seed_place_tags.up.sql

-- Tag the Chrysler Building
INSERT INTO place_tags (place_id, tag_id)
VALUES
  (
    (SELECT id FROM places 
     WHERE name = 'Chrysler Building'),
    (SELECT id FROM tags WHERE slug = 'art-deco')
  ),
  (
    (SELECT id FROM places 
     WHERE name = 'Chrysler Building'),
    (SELECT id FROM tags WHERE slug = 'landmarked')
  ),
  (
    (SELECT id FROM places 
     WHERE name = 'Chrysler Building'),
    (SELECT id FROM tags WHERE slug = 'free-entry')
  );

-- Tag Xi'an Famous Foods
INSERT INTO place_tags (place_id, tag_id)
VALUES
  (
    (SELECT id FROM places 
     WHERE name = 'Xi''an Famous Foods'),
    (SELECT id FROM tags WHERE slug = 'cash-only')
  ),
  (
    (SELECT id FROM places 
     WHERE name = 'Xi''an Famous Foods'),
    (SELECT id FROM tags WHERE slug = 'local-favorite')
  );

-- Tag NYPL Rose Reading Room
INSERT INTO place_tags (place_id, tag_id)
VALUES
  (
    (SELECT id FROM places 
     WHERE name = 'New York Public Library — Rose Main Reading Room'),
    (SELECT id FROM tags WHERE slug = 'free-entry')
  ),
  (
    (SELECT id FROM places 
     WHERE name = 'New York Public Library — Rose Main Reading Room'),
    (SELECT id FROM tags WHERE slug = 'quiet')
  ),
  (
    (SELECT id FROM places 
     WHERE name = 'New York Public Library — Rose Main Reading Room'),
    (SELECT id FROM tags WHERE slug = 'landmarked')
  );

-- Tag Woolworth Building  
INSERT INTO place_tags (place_id, tag_id)
VALUES
  (
    (SELECT id FROM places 
     WHERE name = 'Woolworth Building'),
    (SELECT id FROM tags WHERE slug = 'gothic-revival')
  ),
  (
    (SELECT id FROM places 
     WHERE name = 'Woolworth Building'),
    (SELECT id FROM tags WHERE slug = 'landmarked')
  );