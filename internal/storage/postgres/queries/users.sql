-- name: InsertUser :exec
INSERT INTO users(
    id, username, email, password,
    email_confirmed_at, created_at, updated_at, is_deleted
) VALUES(
    $1, $2, $3, $4, $5,
    $6, $7, $8
);

-- name: SelectUserById :one
SELECT *
FROM users
WHERE id = $1;

-- name: SelectUserByUsername :one
SELECT *
FROM users
WHERE username = $1;

-- name: SelectUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

