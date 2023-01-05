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
    created TIMESTAMP,
    post_tree INT[]
);

CREATE TABLE IF NOT EXISTS votes (
    thread_id INT NOT NULL REFERENCES threads(id) ON DELETE CASCADE,
    nickname CITEXT NOT NULL REFERENCES users(nickname),
    voice INT NOT NULL,
    PRIMARY KEY (thread_id, nickname)
);

CREATE OR REPLACE FUNCTION update_thread_votes_after_insert()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE threads SET votes = votes + NEW.voice WHERE id = NEW.thread_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_vote_trigger
AFTER INSERT ON votes
FOR EACH ROW
EXECUTE PROCEDURE update_thread_votes_after_insert();

CREATE OR REPLACE FUNCTION update_thread_votes_after_update()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE threads SET votes = votes + NEW.voice - OLD.voice WHERE id = NEW.thread_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_vote_trigger
AFTER UPDATE ON votes
FOR EACH ROW
EXECUTE PROCEDURE update_thread_votes_after_update();

CREATE OR REPLACE FUNCTION insert_forum_user()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO forum_user (user_nickname, forum) VALUES (NEW.author, NEW.forum) ON CONFLICT DO NOTHING;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER forum_user_post
AFTER INSERT ON posts
FOR EACH ROW
EXECUTE PROCEDURE insert_forum_user();

CREATE TRIGGER forum_user_thread
AFTER INSERT ON threads
FOR EACH ROW
EXECUTE PROCEDURE insert_forum_user();

CREATE OR REPLACE FUNCTION update_post_tree()
RETURNS TRIGGER AS $$
BEGIN
    NEW.post_tree = (SELECT post_tree FROM posts WHERE id = NEW.parent) || NEW.id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_post_trigger_post
BEFORE INSERT ON posts
FOR EACH ROW
EXECUTE PROCEDURE update_post_tree();

CREATE OR REPLACE FUNCTION update_count_threads()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE forums SET threads = threads + 1 WHERE slug = NEW.forum;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_thread_trigger
AFTER INSERT ON threads
FOR EACH ROW
EXECUTE PROCEDURE update_count_threads();

CREATE OR REPLACE FUNCTION update_count_posts()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE forums SET posts = posts + 1 WHERE slug = NEW.forum;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_post_trigger_forum
AFTER INSERT ON posts
FOR EACH ROW
EXECUTE PROCEDURE update_count_posts();
