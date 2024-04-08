-- name: CreateUser :one
INSERT INTO users(email, is_staff, join_date, last_login)
VALUES (@email, @is_staff, @join_date, @last_login)
RETURNING id;

-- name: GetUserByEmail :one
Select *
FROM users
WHERE id=$1;

-- name: ChangeUserStatus :exec
UPDATE users
SET status=$1
WHERE id=$2;
