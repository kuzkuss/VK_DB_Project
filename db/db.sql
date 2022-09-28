CREATE TABLE forums (
    title VARCHAR(128) NOT NULL,
    user CITEXT NOT NULL REFERENCES users (nickname),
    slug VARCHAR(128) NOT NULL UNIQUE PRIMARY KEY,
    posts INT,
    threads BIGINT
);

CREATE TABLE threads (
    id SERIAL PRIMARY KEY,
    title VARCHAR(128) NOT NULL,
    author CITEXT NOT NULL REFERENCES users (nickname),
    forum VARCHAR(128) REFERENCES forums (slug),
    message TEXT NOT NULL,
    votes INT,
    slug VARCHAR(128),
    created TIMESTAMP
);

CREATE TABLE users (
    nickname CITEXT UNIQUE PRIMARY KEY,
    fullname VARCHAR(128) NOT NULL,
    about TEXT,
    email VARCHAR(128) NOT NULL
);

CREATE TABLE posts (
    id BIGSERIAL,
    parent BIGINT,
    author CITEXT NOT NULL REFERENCES users (nickname),
    message TEXT NOT NULL,
    isEdited BOOLEAN NOT NULL,
    forum VARCHAR(128) REFERENCES forums (slug),
    thread INT REFERENCES threads (id),
    created TIMESTAMP
);

