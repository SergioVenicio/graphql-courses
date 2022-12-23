SELECT
  cat.id
	, cat.name
  , cat.description
FROM categories cat
INNER JOIN courses ON courses.category_id = cat.id 
WHERE courses.id = ?;