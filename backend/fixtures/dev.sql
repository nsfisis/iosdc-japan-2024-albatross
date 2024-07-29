INSERT INTO users
(username, display_name, icon_path, is_admin)
VALUES
    ('a', 'TEST A', NULL, FALSE),
    ('b', 'TEST B', NULL, FALSE),
    ('c', 'TEST C', NULL, TRUE);

INSERT INTO user_auths
(user_id, auth_type)
VALUES
    (1, 'bypass'),
    (2, 'bypass'),
    (3, 'bypass');

INSERT INTO problems
(title, description)
VALUES
    ('TEST problem 1', 'This is TEST problem 1'),
    ('TEST problem 2', 'This is TEST problem 2'),
    ('TEST problem 3', 'This is TEST problem 3');

INSERT INTO games
(state, display_name, duration_seconds, problem_id)
VALUES
    ('waiting_entries', 'TEST game 1', 180, 1),
    ('closed',          'TEST game 2', 180, 2),
    ('finished',        'TEST game 3', 180, 3);

INSERT INTO game_players
(game_id, user_id)
VALUES
    (1, 1),
    (1, 2),
    (2, 1),
    (2, 2),
    (3, 1),
    (3, 2);
