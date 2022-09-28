package storeDB

import (
	"database/sql"

	forums "github.com/kuzkuss/VK_DB_Project/internal/app/models/forumsRepository"
	threads "github.com/kuzkuss/VK_DB_Project/internal/app/models/threadsRepository"
	users "github.com/kuzkuss/VK_DB_Project/internal/app/models/usersRepository"
)

type DataBaseForums struct {
	db *sql.DB
}

type DataBaseUsers struct {
	db *sql.DB
}

type DataBaseThreads struct {
	db *sql.DB
}

func NewDataBaseForums(db *sql.DB) *DataBaseForums {
	return &DataBaseForums {
		db: db,
	}
}

func NewDataBaseUsers(db *sql.DB) *DataBaseUsers {
	return &DataBaseUsers {
		db: db,
	}
}

func NewDataBaseThreads(db *sql.DB) *DataBaseThreads {
	return &DataBaseThreads {
		db: db,
	}
}

func (dbForums *DataBaseForums) CreateForum(f forums.Forum) (*forums.Forum, error) {
	row, err := dbForums.db.Query("INSERT INTO forums VALUES ($1, $2, $3, $4, $5) RETURNING title, user, slug, posts, threads", 
	f.Title, f.User, f.Slug, f.Posts, f.Threads)
	if err != nil {
		return nil, err
  	}
	
	forum := forums.Forum{}

	row.Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)

	return &forum, nil
}

func (dbThreads *DataBaseThreads) CreateThread(t threads.Thread) (*threads.Thread, error) {
	row, err := dbThreads.db.Query("INSERT INTO threads (title, author, forum, message, votes, slug, created) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, title, author, forum, message, votes, slug, created", 
	t.Id, t.Title, t.Author, t.Forum, t.Message, t.Votes, t.Slug, t.Created)

	if err != nil {
		return nil, err
  	}
	
	thread := threads.Thread{}

	row.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message,
		&thread.Votes, &thread.Slug, &thread.Created)

	_, err = dbThreads.db.Query("UPDATE forums SET threads = threads + 1 WHERE forum.slug = " + thread.Forum)
	if err != nil {
		return nil, err
	}

	return &thread, nil
}

func (dbForums *DataBaseForums) SelectForum(slug string) (*forums.Forum, error) {
	row, err := dbForums.db.Query("SELECT * FROM forums WHERE slug=" + slug)
	if err != nil {
		return nil, err
  	}
	
	forum := forums.Forum{}

	row.Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)

	return &forum, nil
}

func (dbThreads *DataBaseThreads) SelectThread(slug string) (*threads.Thread, error) {
	row, err := dbThreads.db.Query("SELECT * FROM threads WHERE slug=" + slug)
	if err != nil {
		return nil, err
  	}
	
	thread := threads.Thread{}

	row.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message,
		&thread.Votes, &thread.Slug, &thread.Created)

	return &thread, nil
}

func (dbUsers *DataBaseUsers) SelectUser(name string) (*users.User, error) {
	row, err := dbUsers.db.Query("SELECT * FROM users WHERE nickname=" + name)
	if err != nil {
		return nil, err
  	}
	
	user := users.User{}

	row.Scan(&user.Nickname, user.Fullname, user.About, user.Email)

	return &user, nil
}

