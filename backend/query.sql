-- name: GetUserByID :one
SELECT * FROM users
WHERE users.user_id = $1
LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users;

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

-- name: UpdateGameState :exec
UPDATE games
SET state = $2
WHERE game_id = $1;

-- name: UpdateGameStartedAt :exec
UPDATE games
SET started_at = $2
WHERE game_id = $1;

-- name: GetGameByID :one
SELECT * FROM games
LEFT JOIN problems ON games.problem_id = problems.problem_id
WHERE games.game_id = $1
LIMIT 1;

-- name: ListGamePlayers :many
SELECT * FROM game_players
LEFT JOIN users ON game_players.user_id = users.user_id
WHERE game_players.game_id = $1;

-- name: UpdateGame :exec
UPDATE games
SET
    game_type = $2,
    state = $3,
    display_name = $4,
    duration_seconds = $5,
    started_at = $6,
    problem_id = $7
WHERE game_id = $1;
