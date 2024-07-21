-- name: GetUserAuthFromUsername :one
SELECT * FROM users
JOIN user_auths ON users.user_id = user_auths.user_id
WHERE users.username = $1;
