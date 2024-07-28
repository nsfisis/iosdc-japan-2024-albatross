CREATE TABLE users (
    user_id      SERIAL      PRIMARY KEY,
    username     VARCHAR(64) NOT NULL UNIQUE,
    display_name VARCHAR(64) NOT NULL,
    icon_path    VARCHAR(255),
    is_admin     BOOLEAN     NOT NULL,
    created_at   TIMESTAMP   NOT NULL DEFAULT NOW()
);

CREATE TABLE user_auths (
    user_auth_id  SERIAL      PRIMARY KEY,
    user_id       INT         NOT NULL,
    auth_type     VARCHAR(16) NOT NULL,
    password_hash VARCHAR(256),
    CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES users(user_id)
);

CREATE TABLE games (
    game_id          SERIAL       PRIMARY KEY,
    state            VARCHAR(32)  NOT NULL,
    display_name     VARCHAR(255) NOT NULL,
    duration_seconds INT          NOT NULL,
    created_at       TIMESTAMP    NOT NULL DEFAULT NOW(),
    started_at       TIMESTAMP,
    problem_id       INT,
    CONSTRAINT fk_problem_id FOREIGN KEY(problem_id) REFERENCES problems(problem_id)
);

CREATE TABLE game_players (
    game_id INT NOT NULL,
    user_id INT NOT NULL,
    PRIMARY KEY (game_id, user_id),
    CONSTRAINT fk_game_id FOREIGN KEY(game_id) REFERENCES games(game_id),
    CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES users(user_id)
);

CREATE TABLE problems (
    problem_id  SERIAL       PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL
);

CREATE TABLE testcases (
    testcase_id SERIAL PRIMARY KEY,
    problem_id  INT    NOT NULL,
    stdin       TEXT   NOT NULL,
    stdout      TEXT   NOT NULL,
    CONSTRAINT fk_problem_id FOREIGN KEY(problem_id) REFERENCES problems(problem_id)
);
