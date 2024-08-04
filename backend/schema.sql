CREATE TABLE users (
    user_id      SERIAL      PRIMARY KEY,
    username     VARCHAR(64) NOT NULL UNIQUE,
    display_name VARCHAR(64) NOT NULL,
    icon_path    VARCHAR(255),
    is_admin     BOOLEAN     NOT NULL,
    created_at   TIMESTAMP   NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_users_username ON users(username);

CREATE TABLE user_auths (
    user_auth_id  SERIAL      PRIMARY KEY,
    user_id       INT         NOT NULL UNIQUE,
    auth_type     VARCHAR(16) NOT NULL,
    password_hash VARCHAR(256),
    CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES users(user_id)
);
CREATE INDEX idx_user_auths_user_id ON user_auths(user_id);

CREATE TABLE problems (
    problem_id  SERIAL       PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL
);

CREATE TABLE games (
    game_id          SERIAL       PRIMARY KEY,
    game_type        VARCHAR(16)  NOT NULL,
    state            VARCHAR(32)  NOT NULL,
    display_name     VARCHAR(255) NOT NULL,
    duration_seconds INT          NOT NULL,
    created_at       TIMESTAMP    NOT NULL DEFAULT NOW(),
    started_at       TIMESTAMP,
    problem_id       INT,
    CONSTRAINT fk_problem_id FOREIGN KEY(problem_id) REFERENCES problems(problem_id)
);
CREATE INDEX idx_games_problem_id ON games(problem_id);

CREATE TABLE game_players (
    game_id INT NOT NULL,
    user_id INT NOT NULL,
    PRIMARY KEY (game_id, user_id),
    CONSTRAINT fk_game_id FOREIGN KEY(game_id) REFERENCES games(game_id),
    CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES users(user_id)
);

CREATE TABLE testcases (
    testcase_id SERIAL PRIMARY KEY,
    problem_id  INT    NOT NULL,
    stdin       TEXT   NOT NULL,
    stdout      TEXT   NOT NULL,
    CONSTRAINT fk_problem_id FOREIGN KEY(problem_id) REFERENCES problems(problem_id)
);
CREATE INDEX idx_testcases_problem_id ON testcases(problem_id);

CREATE TABLE submissions (
    submission_id SERIAL    PRIMARY KEY,
    game_id       INT       NOT NULL,
    user_id       INT       NOT NULL,
    code          TEXT      NOT NULL,
    code_size     INT       NOT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_game_id FOREIGN KEY(game_id) REFERENCES games(game_id),
    CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES users(user_id)
);

CREATE TABLE testcase_executions (
    testcase_execution_id SERIAL      PRIMARY KEY,
    submission_id         INT         NOT NULL,
    testcase_id           INT,
    status                VARCHAR(16) NOT NULL,
    stdout                TEXT        NOT NULL,
    stderr                TEXT        NOT NULL,
    CONSTRAINT fk_submission_id FOREIGN KEY(submission_id) REFERENCES submissions(submission_id),
    CONSTRAINT fk_testcase_id FOREIGN KEY(testcase_id) REFERENCES testcases(testcase_id)
);
CREATE INDEX idx_testcase_executions_submission_id ON testcase_executions(submission_id);
