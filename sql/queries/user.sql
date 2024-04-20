-- name: AllStaff :many
SELECT *
FROM users
WHERE is_staff = TRUE
ORDER BY join_date DESC
LIMIT $1 OFFSET $2;

-- name: AllStaffCount :one
SELECT COUNT(*)
FROM users
WHERE is_staff = TRUE;

-- name: CreateUser :one
INSERT INTO users(name, email, phone_number, avatar, is_staff, join_date, last_login)
VALUES (@name, @email, @phone_number, @avatar, @is_staff, NOW(), NOW())
RETURNING *;

-- name: GetUserByEmail :one
Select *
FROM users
WHERE email = @email;

-- name: GetUSerById :one
SELECT *
FROM users
WHERE id = @id;

-- name: UpdateUser :one
UPDATE users
SET name         = $2,
    email        = $3,
    phone_number = $4,
    status       = $5
WHERE id = $1
RETURNING *;
