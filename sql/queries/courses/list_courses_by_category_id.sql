SELECT
  id
  , name
  , description
FROM courses
WHERE category_id = ?
ORDER BY name;