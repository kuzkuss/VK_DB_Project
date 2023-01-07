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
    threads INT DEFAULT 0
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
    parent INT,
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




-- Indexes

-- CREATE INDEX IF NOT EXISTS index_forums_user_nickname ON forums (user_nickname);
-- CREATE INDEX IF NOT EXISTS index_threads_author ON threads (author);
-- CREATE INDEX IF NOT EXISTS index_threads_forum ON threads using hash (forum);
-- CREATE INDEX IF NOT EXISTS index_forum_user_user_nickname ON forum_user (user_nickname);
-- CREATE INDEX IF NOT EXISTS index_forum_user_forum ON forum_user (forum);
-- CREATE INDEX IF NOT EXISTS index_posts_author ON posts using hash (author);
-- CREATE INDEX IF NOT EXISTS index_posts_forum ON posts using hash (forum);
-- CREATE INDEX IF NOT EXISTS index_posts_thread ON posts (thread);
-- CREATE INDEX IF NOT EXISTS index_votes_thread_id ON votes (thread_id);
-- CREATE INDEX IF NOT EXISTS index_votes_nickname ON votes (nickname);



-- CREATE INDEX IF NOT EXISTS  index_users_email ON users (email);

-- -- Threads
-- CREATE INDEX IF NOT EXISTS index_slug_thread on threads using hash (slug);

-- CREATE INDEX IF NOT EXISTS index_thread_forum_created ON threads (forum, created);

-- -- Post
-- CREATE INDEX IF NOT EXISTS index_parent_post on posts (parent);
-- CREATE INDEX IF NOT EXISTS index_thread_post_tree_post on posts (thread, post_tree);
-- CREATE INDEX IF NOT EXISTS index_first_parent_post on posts ((post_tree[1]), post_tree);

-- -- Vote
-- CREATE INDEX IF NOT EXISTS index_search_user_vote ON votes (nickname, thread_id, voice);

-- VACUUM ANALYZE;


-- Indexes change


CREATE INDEX IF NOT EXISTS index_forums_user_nickname ON forums (user_nickname);
CREATE INDEX IF NOT EXISTS index_threads_author ON threads (author);
CREATE INDEX IF NOT EXISTS index_threads_forum ON threads (forum);
CREATE INDEX IF NOT EXISTS index_posts_author ON posts (author);
CREATE INDEX IF NOT EXISTS index_posts_forum ON posts (forum);
CREATE INDEX IF NOT EXISTS index_posts_thread ON posts (thread);
CREATE INDEX IF NOT EXISTS index_votes_thread_id ON votes (thread_id);
CREATE INDEX IF NOT EXISTS index_votes_nickname ON votes (nickname);

CREATE INDEX IF NOT EXISTS index_forum_user_forum ON forum_user (forum, user_nickname);



CREATE INDEX IF NOT EXISTS  index_users_email ON users (email);

-- Threads
CREATE INDEX IF NOT EXISTS index_slug_thread on threads (slug);

CREATE INDEX IF NOT EXISTS index_thread_forum_created ON threads (forum, created);

-- Post
CREATE INDEX IF NOT EXISTS index_parent_post on posts (parent);
CREATE INDEX IF NOT EXISTS index_thread_post_tree_post on posts (thread, post_tree);
CREATE INDEX IF NOT EXISTS index_first_parent_post on posts ((post_tree[1]), post_tree);

-- Vote
CREATE INDEX IF NOT EXISTS index_search_user_vote ON votes (nickname, thread_id, voice);

VACUUM ANALYZE;

