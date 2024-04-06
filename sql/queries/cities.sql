-- name: CreateCity :one
INSERT INTO cities
VALUES ($1, $2)
RETURNING id;

-- name: ListAllCities :many
SELECT *
FROM cities
LIMIT $1 OFFSET $2;

-- name: ListActiveCity :many
SELECT *
FROM cities
WHERE is_active=true
LIMIT $1 OFFSET $2;
