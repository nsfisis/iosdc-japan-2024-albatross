-- name: GetUserByID :one
SELECT * FROM users
WHERE users.user_id = $1
LIMIT 1;

-- name: GetUserIDByUsername :one
SELECT user_id FROM users
WHERE users.username = $1
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (username, display_name, is_admin)
VALUES ($1, $1, false)
RETURNING user_id;

-- name: UpdateUserIconPath :exec
UPDATE users
SET icon_path = $2
WHERE user_id = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY users.user_id;

-- name: GetUserAuthByUsername :one
SELECT * FROM users
JOIN user_auths ON users.user_id = user_auths.user_id
WHERE users.username = $1
LIMIT 1;

-- name: CreateUserAuth :exec
INSERT INTO user_auths (user_id, auth_type)
VALUES ($1, $2);

-- name: IsRegistrationTokenValid :one
SELECT EXISTS (
    SELECT 1 FROM registration_tokens
    WHERE token = $1
);

-- name: ListGames :many
SELECT * FROM games
JOIN problems ON games.problem_id = problems.problem_id
ORDER BY games.game_id;

-- name: ListGamesForPlayer :many
SELECT * FROM games
JOIN problems ON games.problem_id = problems.problem_id
JOIN game_players ON games.game_id = game_players.game_id
WHERE game_players.user_id = $1
ORDER BY games.game_id;

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
JOIN problems ON games.problem_id = problems.problem_id
WHERE games.game_id = $1
LIMIT 1;

-- name: ListGamePlayers :many
SELECT * FROM game_players
JOIN users ON game_players.user_id = users.user_id
WHERE game_players.game_id = $1
ORDER BY game_players.user_id;

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

-- name: CreateSubmission :one
INSERT INTO submissions (game_id, user_id, code, code_size, code_hash)
VALUES ($1, $2, $3, $4, $5)
RETURNING submission_id;

-- name: GetSubmissionCodeSizeByID :one
SELECT code_size FROM submissions
WHERE submission_id = $1
LIMIT 1;

-- name: ListTestcasesByGameID :many
SELECT * FROM testcases
WHERE testcases.problem_id = (SELECT problem_id FROM games WHERE game_id = $1)
ORDER BY testcases.testcase_id;

-- name: ListTestcaseIDsByGameID :many
SELECT testcases.testcase_id FROM testcases
WHERE testcases.problem_id = (SELECT problem_id FROM games WHERE game_id = $1)
ORDER BY testcases.testcase_id;

-- name: CreateSubmissionResult :exec
INSERT INTO submission_results (submission_id, status, stdout, stderr)
VALUES ($1, $2, $3, $4);

-- name: CreateTestcaseResult :exec
INSERT INTO testcase_results (submission_id, testcase_id, status, stdout, stderr)
VALUES ($1, $2, $3, $4, $5);

-- name: AggregateTestcaseResults :one
SELECT
    CASE
        WHEN COUNT(CASE WHEN r.status IS NULL            THEN 1 END) > 0 THEN 'running'
        WHEN COUNT(CASE WHEN r.status = 'internal_error' THEN 1 END) > 0 THEN 'internal_error'
        WHEN COUNT(CASE WHEN r.status = 'timeout'        THEN 1 END) > 0 THEN 'timeout'
        WHEN COUNT(CASE WHEN r.status = 'runtime_error'  THEN 1 END) > 0 THEN 'runtime_error'
        WHEN COUNT(CASE WHEN r.status = 'wrong_answer'   THEN 1 END) > 0 THEN 'wrong_answer'
        ELSE 'success'
    END AS status
FROM testcases
LEFT JOIN testcase_results AS r ON testcases.testcase_id = r.testcase_id
WHERE r.submission_id = $1;

-- A submission is accepted if its result is "success" and no testcase failed.
-- Note that testcase_results only has rows for failed testcases, and that
-- submission_results is written as soon as the first testcase finishes.

-- name: GetRanking :many
SELECT
    best.submission_id,
    best.code,
    best.code_size,
    best.created_at,
    users.user_id,
    users.username,
    users.display_name,
    users.icon_path,
    users.is_admin
FROM (
    SELECT DISTINCT ON (submissions.user_id)
        submissions.submission_id,
        submissions.user_id,
        submissions.code,
        submissions.code_size,
        submissions.created_at
    FROM submissions
    JOIN submission_results ON submissions.submission_id = submission_results.submission_id
    WHERE submissions.game_id = $1
        AND submission_results.status = 'success'
        AND NOT EXISTS (
            SELECT 1 FROM testcase_results
            WHERE testcase_results.submission_id = submissions.submission_id
        )
    ORDER BY submissions.user_id, submissions.code_size, submissions.created_at, submissions.submission_id
) AS best
JOIN users ON best.user_id = users.user_id
ORDER BY best.code_size, best.created_at, best.submission_id;

-- name: ListSubmissionSummaries :many
SELECT
    submissions.submission_id,
    submissions.user_id,
    submissions.code_size,
    submissions.created_at,
    submission_results.status,
    (
        SELECT COUNT(*) FROM testcase_results
        WHERE testcase_results.submission_id = submissions.submission_id
    ) AS failed_testcase_count
FROM submissions
LEFT JOIN submission_results ON submissions.submission_id = submission_results.submission_id
WHERE submissions.game_id = $1
ORDER BY submissions.submission_id;
