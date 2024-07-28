-- name: GetUserById :one
SELECT * FROM users
WHERE users.user_id = $1
LIMIT 1;

-- name: GetUserAuthByUsername :one
SELECT * FROM users
JOIN user_auths ON users.user_id = user_auths.user_id
WHERE users.username = $1
LIMIT 1;

-- name: ListGames :many
SELECT * FROM games
LEFT JOIN problems ON games.problem_id = problems.problem_id;

-- name: ListGamesForPlayer :many
SELECT * FROM games
LEFT JOIN problems ON games.problem_id = problems.problem_id
JOIN game_players ON games.game_id = game_players.game_id
WHERE game_players.user_id = $1;
