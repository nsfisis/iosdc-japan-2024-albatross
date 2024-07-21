CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(64) NOT NULL UNIQUE,
    display_username VARCHAR(64) NOT NULL,
    icon_url VARCHAR(255),
    is_admin BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE user_auths (
    user_auth_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    auth_type VARCHAR(16) NOT NULL,
    password_hash VARCHAR(256),
    CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES users(user_id)
);

CREATE TABLE games (
    game_id SERIAL PRIMARY KEY,
    type VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    state VARCHAR(255) NOT NULL
);
