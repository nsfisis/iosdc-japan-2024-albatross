INSERT INTO users
(username, display_name, icon_path, is_admin)
VALUES
    ('a', 'TEST A', NULL, FALSE),
    ('b', 'TEST B', NULL, FALSE),
    ('c', 'TEST C', NULL, TRUE);

INSERT INTO user_auths
(user_id, auth_type, password_hash)
VALUES
    (1, 'password', '$2a$10$5FzjoitnZSFFpIPHEqmnXOQkSKWTLwpR.gqPy50iFg5itOZcqARHq'),
    (2, 'password', '$2a$10$4Wl1M4jQs.GwkB4oT32KvuMQtF.EdqKuOc8z8KKOupnuMJRAVk32W'),
    (3, 'password', '$2a$10$F/TePpu1pyJRWgn0e6A14.VL9D/17sRxT/2DyZ2Oi4Eg/lR6n7JcK');

INSERT INTO problems
(title, description)
VALUES
    ('TEST problem 1', 'This is TEST problem 1'),
    ('TEST problem 2', 'This is TEST problem 2'),
    ('TEST problem 3', 'This is TEST problem 3'),
    ('TEST problem 4', 'This is TEST problem 4'),
    ('TEST problem 5', 'This is TEST problem 5'),
    ('TEST problem 6', 'This is TEST problem 6'),
    ('TEST problem 7', 'This is TEST problem 7');

INSERT INTO games
(game_type, state, display_name, duration_seconds, problem_id)
VALUES
    ('1v1',         'waiting_entries', 'TEST game 1', 180, 1),
    ('1v1',         'closed',          'TEST game 2', 180, 2),
    ('1v1',         'finished',        'TEST game 3', 180, 3),
    ('multiplayer', 'waiting_start',   'TEST game 4', 180, 4),
    ('multiplayer', 'closed',          'TEST game 5', 180, 5),
    ('multiplayer', 'finished',        'TEST game 6', 180, 6),
    ('multiplayer', 'waiting_entries', 'TEST game 7', 180, 7);

INSERT INTO game_players
(game_id, user_id)
VALUES
    (1, 1),
    (1, 2),
    (2, 1),
    (2, 2),
    (3, 1),
    (3, 2),
    (4, 1),
    (4, 2),
    (5, 1),
    (5, 2),
    (6, 1),
    (6, 2),
    (7, 1);

INSERT INTO testcases
(problem_id, stdin, stdout)
VALUES
    (1, '', '42'),
    (4, '', '42'),
    (7, '', '42');
