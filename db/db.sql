CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
    nickname CITEXT PRIMARY KEY,
    fullname VARCHAR(128) NOT NULL,
    about TEXT,
    email CITEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS forums (
    title VARCHAR(128) NOT NULL,
    user_nickname CITEXT NOT NULL REFERENCES users(nickname),
    slug CITEXT PRIMARY KEY,
    posts INT DEFAULT 0,
    threads INT DEFAULT 0 -- был bigint
);

CREATE TABLE IF NOT EXISTS threads (
    id SERIAL PRIMARY KEY,
    title VARCHAR(128) NOT NULL,
    author CITEXT NOT NULL REFERENCES users(nickname),
    forum CITEXT NOT NULL REFERENCES forums(slug) ON DELETE CASCADE,
    message TEXT NOT NULL,
    votes INT DEFAULT 0,
    slug CITEXT UNIQUE,
    created TIMESTAMP
);

CREATE TABLE IF NOT EXISTS forum_user (
    user_nickname CITEXT NOT NULL REFERENCES users(nickname) ON DELETE CASCADE,
    forum CITEXT NOT NULL REFERENCES forums(slug) ON DELETE CASCADE,
    PRIMARY KEY (user_nickname, forum)
);

CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    parent INT, -- был bigint
    author CITEXT NOT NULL REFERENCES users(nickname),
    message TEXT NOT NULL,
    is_edited BOOLEAN NOT NULL,
    forum CITEXT REFERENCES forums(slug) ON DELETE CASCADE,
    thread INT REFERENCES threads(id) ON DELETE CASCADE,
    created TIMESTAMP
);

CREATE TABLE IF NOT EXISTS votes (
    thread_id INT NOT NULL REFERENCES threads(id) ON DELETE CASCADE,
    nickname CITEXT NOT NULL REFERENCES users(nickname),
    voice INT NOT NULL,
    PRIMARY KEY (thread_id, nickname)
)



