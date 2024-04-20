-- name: AllCities :many
SELECT *
FROM cities
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: CitiesCount :one
SELECT COUNT(*)
FROM cities;

-- name: FilterCities :many
SELECT *
FROM cities
WHERE name_en ILIKE @query
   or name_ar ILIKE @query
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: ActiveCities :many
SELECT *
FROM cities
WHERE is_active = TRUE
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: ActiveCitiesCount :one
SELECT COUNT(*)
FROM cities
WHERE is_active = TRUE;

-- name: CreateCity :one
INSERT INTO cities(name_en, name_ar, is_active)
VALUES (@name_en, @name_ar, @is_active)
RETURNING id;

-- name: UpdateCity :one
UPDATE cities
SET name_en=$1,
    name_ar=$2,
    is_active=$3
WHERE id = $4
RETURNING *;

-- name: DeleteCity :exec
DELETE
FROM cities
WHERE id = $1;
