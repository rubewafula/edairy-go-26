-- Dev seed data: collection centers for routes
-- Every route must have at least one center; milk is always recorded against a center.
-- Run against dev_dairy_erp AFTER seed-milk-journal-import.sql (which names DEV-R01..DEV-R10):
--   mysql -uedairy -p dev_dairy_erp < deploy/dev/seed-route-centers.sql
-- Safe to re-run (idempotent).

-- 1. Named centers for the dev routes (two per route)
INSERT INTO route_centers (route_id, center, created_by, updated_by)
SELECT r.id, seed.center, 1, 1
FROM (
  SELECT 'DEV-R01' AS route_code, 'Dev R01 Center A' AS center UNION ALL
  SELECT 'DEV-R01', 'Dev R01 Center B' UNION ALL
  SELECT 'DEV-R02', 'Dev R02 Center A' UNION ALL
  SELECT 'DEV-R02', 'Dev R02 Center B' UNION ALL
  SELECT 'DEV-R03', 'Dev R03 Center A' UNION ALL
  SELECT 'DEV-R03', 'Dev R03 Center B' UNION ALL
  SELECT 'DEV-R04', 'Dev R04 Center A' UNION ALL
  SELECT 'DEV-R04', 'Dev R04 Center B' UNION ALL
  SELECT 'DEV-R05', 'Dev R05 Center A' UNION ALL
  SELECT 'DEV-R05', 'Dev R05 Center B' UNION ALL
  SELECT 'DEV-R06', 'Dev R06 Center A' UNION ALL
  SELECT 'DEV-R06', 'Dev R06 Center B' UNION ALL
  SELECT 'DEV-R07', 'Dev R07 Center A' UNION ALL
  SELECT 'DEV-R07', 'Dev R07 Center B' UNION ALL
  SELECT 'DEV-R08', 'Dev R08 Center A' UNION ALL
  SELECT 'DEV-R08', 'Dev R08 Center B' UNION ALL
  SELECT 'DEV-R09', 'Dev R09 Center A' UNION ALL
  SELECT 'DEV-R09', 'Dev R09 Center B' UNION ALL
  SELECT 'DEV-R10', 'Dev R10 Center A' UNION ALL
  SELECT 'DEV-R10', 'Dev R10 Center B'
) AS seed
JOIN routes r ON r.route_code = seed.route_code AND r.deleted_at IS NULL
WHERE NOT EXISTS (
  SELECT 1 FROM route_centers rc
  WHERE rc.route_id = r.id AND rc.center = seed.center AND rc.deleted_at IS NULL
);

-- 2. Safety net: any other route still without a center gets a default one
INSERT INTO route_centers (route_id, center, created_by, updated_by)
SELECT r.id,
       LEFT(CONCAT(COALESCE(NULLIF(r.route_name, ''), NULLIF(r.route_code, ''), CONCAT('Route ', r.id)), ' Center'), 96),
       1, 1
FROM routes r
WHERE r.deleted_at IS NULL
  AND NOT EXISTS (
    SELECT 1 FROM route_centers rc WHERE rc.route_id = r.id AND rc.deleted_at IS NULL
  );

-- 3. Backfill existing milk journal entries that have no center:
--    assign the first center of the journal's route
UPDATE milk_journal_entries e
JOIN milk_journals j ON j.id = e.milk_journal_id
JOIN (
  SELECT route_id, MIN(id) AS center_id
  FROM route_centers
  WHERE deleted_at IS NULL
  GROUP BY route_id
) c ON c.route_id = j.route_id
SET e.route_center_id = c.center_id
WHERE (e.route_center_id IS NULL OR e.route_center_id = 0);

-- 4. Check: should return 0 rows
SELECT r.id, r.route_code, r.route_name
FROM routes r
WHERE r.deleted_at IS NULL
  AND NOT EXISTS (SELECT 1 FROM route_centers rc WHERE rc.route_id = r.id AND rc.deleted_at IS NULL);
