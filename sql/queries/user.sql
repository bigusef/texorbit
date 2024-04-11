-- name: AllStaffUser :many
SELECT *
FROM users
WHERE is_staff = TRUE
ORDER BY join_date DESC
LIMIT $1 OFFSET $2;

-- name: StaffCount :one
SELECT COUNT(*)
FROM users
WHERE is_staff = TRUE;

-- name: GetUserByEmail :one
Select *
FROM users
WHERE email = @email;

-- name: GetUSerById :one
SELECT *
FROM users
WHERE id = @id;

-- name: CreateUser :one
INSERT INTO users(name, email, avatar, is_staff, join_date, last_login)
VALUES (@name, @email, @avatar, @is_staff, NOW(), NOW())
RETURNING *;



