-- name: GetUserById :one
SELECT * FROM users
WHERE users.user_id = $1
LIMIT 1;

-- name: GetUserAuthByUsername :one
SELECT * FROM users
JOIN user_auths ON users.user_id = user_auths.user_id
WHERE users.username = $1
LIMIT 1;
