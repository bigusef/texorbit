-- name: CreateCity :one
INSERT INTO cities(name_en, name_ar, is_active)
VALUES (@name_en, @name_ar, @is_active)
RETURNING id;

-- name: CitiesCount :one
SELECT COUNT(*)
FROM cities;

-- name: ListAllCities :many
SELECT *
FROM cities
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: ListActiveCity :many
SELECT *
FROM cities
WHERE is_active = true
ORDER BY id
LIMIT $1 OFFSET $2;

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
