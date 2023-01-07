package repository

import (
	// "context"

	"github.com/kuzkuss/VK_DB_Project/app/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	// "github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryI interface {
	CreatePosts(posts []*models.Post) (error)
	UpdatePost(post *models.Post) (error)
	SelectPostById(id uint64) (*models.Post, error)
	SelectThreadPosts(id uint64, limit int, since int, desc bool, sort string) ([]*models.Post, error)
}

type dataBase struct {
	db *gorm.DB
	// pool *pgxpool.Pool
}

func New(db *gorm.DB) RepositoryI {
	return &dataBase{
		db: db,
		// pool: pool,
	}
}

func (dbPost *dataBase) CreatePosts(posts []*models.Post) (error) {
	tx := dbPost.db.Create(&posts)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table posts)")
	}

	return nil
}

func (dbPost *dataBase) UpdatePost(post *models.Post) (error) {
	tx := dbPost.db.Model(post).Clauses(clause.Returning{}).Updates(models.Post{Message:post.Message, IsEdited: true})
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table posts)")
	}

	return nil
}

func (dbPost *dataBase) SelectPostById(id uint64) (*models.Post, error) {
	post := models.Post{}

	tx := dbPost.db.Where("id = ?", id).Take(&post)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table posts)")
	}

	return &post, nil
}

func (dbPost *dataBase) SelectThreadPosts(id uint64, limit int, since int, desc bool, sort string) ([]*models.Post, error) {
	posts := make([]*models.Post, 0, 10)

	if sort == "flat" {
		if desc {
			if since != 0 {
				tx := dbPost.db.Limit(limit).Where("thread = ? AND id < ?", id, since).
				Order("id desc").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			} else {
				tx := dbPost.db.Limit(limit).Where("thread = ?", id).Order("id desc").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			}
		} else {
			if since != 0 {
				tx := dbPost.db.Limit(limit).Where("thread = ? AND id > ?", id, since).
				Order("id").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			} else {
				tx := dbPost.db.Limit(limit).Where("thread = ?", id).Order("id").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			}
		}
	} else if sort == "tree" {
		if desc {
			if since != 0 {
				tx := dbPost.db.Limit(limit).Where("thread = ? AND post_tree < (?)", id,
				dbPost.db.Table("posts").Select("post_tree").Where("id = ?", since)).
				Order("post_tree desc").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			} else {
				tx := dbPost.db.Limit(limit).Where("thread = ?", id).Order("post_tree desc").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			}
		} else {
			if since != 0 {
				tx := dbPost.db.Limit(limit).Where("thread = ? AND post_tree > (?)", id,
				dbPost.db.Table("posts").Select("post_tree").Where("id = ?", since)).
				Order("post_tree").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			} else {
				tx := dbPost.db.Limit(limit).Where("thread = ?", id).Order("post_tree").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			}
		}
	} else if sort == "parent_tree" {
		if desc {
			if since != 0 {
				tx := dbPost.db.Where("post_tree[1] IN (?)", dbPost.db.
				Table("posts").Limit(limit).Where("parent = 0 AND thread = ? AND id < (?)", id, 
				dbPost.db.Table("posts").Select("post_tree[1]").Where("id = ?", since)).
				Order("id desc").Select("id")).Order("post_tree[1] desc, post_tree").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			} else {
				tx := dbPost.db.Where("post_tree[1] IN (?)", dbPost.db.
				Table("posts").Limit(limit).Where("parent = 0 AND thread = ?", id).
				Order("id desc").Select("id")).Order("post_tree[1] desc, post_tree").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			}
		} else {
			if since != 0 {
				tx := dbPost.db.Where("post_tree[1] IN (?)", dbPost.db.
				Table("posts").Limit(limit).Where("parent = 0 AND thread = ? AND id > (?)", id,
				dbPost.db.Table("posts").Select("post_tree[1]").Where("id = ?", since)).
				Order("id").Select("id")).Order("post_tree").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			} else {
				tx := dbPost.db.Where("post_tree[1] IN (?)", dbPost.db.
				Table("posts").Limit(limit).Where("parent = 0 AND thread = ?", id).
				Order("id").Select("id")).Order("post_tree").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			}
		}
	}

	return posts, nil
}

// func (dbPost *dataBase) SelectThreadPosts(id uint64, limit int, since int, desc bool, sort string) ([]*models.Post, error) {
// 	posts := make([]*models.Post, 0, 10)

// 	if sort == "flat" {
// 		if desc {
// 			if since != 0 {
// 				query := `SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts
// 				  WHERE thread = $1 AND id < $2
// 				  ORDER BY id DESC
// 				  LIMIT $3`
// 				rows, err := dbPost.pool.Query(context.Background(), query, id, since, limit)
// 				if err != nil {
// 					return nil, err
// 				}

// 				for rows.Next() {
// 					element := models.Post{}
// 					if err := rows.Scan(&element.Id, &element.Parent, &element.Author, 
// 						&element.Message, &element.IsEdited, &element.Forum, &element.Thread, 
// 						&element.Created); err != nil {
// 						return nil, err
// 					}
// 					posts = append(posts, &element)
// 				}
// 			} else {
// 				query := `SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts
// 				  WHERE thread = $1
// 				  ORDER BY id DESC
// 				  LIMIT $2`

// 				rows, err := dbPost.pool.Query(context.Background(), query, id, limit)
// 				if err != nil {
// 					return nil, err
// 				}

// 				for rows.Next() {
// 					element := models.Post{}
// 					if err := rows.Scan(&element.Id, &element.Parent, &element.Author, 
// 						&element.Message, &element.IsEdited, &element.Forum, &element.Thread, 
// 						&element.Created); err != nil {
// 						return nil, err
// 					}
// 					posts = append(posts, &element)
// 				}
// 			}
// 		} else {
// 			if since != 0 {
// 				query := `SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts
// 				  WHERE thread = $1 AND id > $2
// 				  ORDER BY id
// 				  LIMIT $3`

// 				rows, err := dbPost.pool.Query(context.Background(), query, id, since, limit)
// 				if err != nil {
// 					return nil, err
// 				}

// 				for rows.Next() {
// 					element := models.Post{}
// 					if err := rows.Scan(&element.Id, &element.Parent, &element.Author, 
// 						&element.Message, &element.IsEdited, &element.Forum, &element.Thread, 
// 						&element.Created); err != nil {
// 						return nil, err
// 					}
// 					posts = append(posts, &element)
// 				}
// 			} else {
// 				query := `SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts
// 				  WHERE thread = $1
// 				  ORDER BY id
// 				  LIMIT $2`

// 				rows, err := dbPost.pool.Query(context.Background(), query, id, limit)
// 				if err != nil {
// 					return nil, err
// 				}

// 				for rows.Next() {
// 					element := models.Post{}
// 					if err := rows.Scan(&element.Id, &element.Parent, &element.Author, 
// 						&element.Message, &element.IsEdited, &element.Forum, &element.Thread, 
// 						&element.Created); err != nil {
// 						return nil, err
// 					}
// 					posts = append(posts, &element)
// 				}
// 			}
// 		}
// 	} else if sort == "tree" {
// 		if desc {
// 			if since != 0 {
// 				query := `SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts
// 				  WHERE thread = $1 and post_tree < (select post_tree from posts where id = $2)
// 				  ORDER BY post_tree DESC
// 				  LIMIT $3;`

// 				rows, err := dbPost.pool.Query(context.Background(), query, id, since, limit)
// 				if err != nil {
// 					return nil, err
// 				}

// 				for rows.Next() {
// 					element := models.Post{}
// 					if err := rows.Scan(&element.Id, &element.Parent, &element.Author, 
// 						&element.Message, &element.IsEdited, &element.Forum, &element.Thread, 
// 						&element.Created); err != nil {
// 						return nil, err
// 					}
// 					posts = append(posts, &element)
// 				}
// 			} else {
// 				query := `SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts
// 				  WHERE thread = $1
// 				  ORDER BY post_tree DESC
// 				  LIMIT $2;`

// 				rows, err := dbPost.pool.Query(context.Background(), query, id, limit)
// 				if err != nil {
// 					return nil, err
// 				}

// 				for rows.Next() {
// 					element := models.Post{}
// 					if err := rows.Scan(&element.Id, &element.Parent, &element.Author, 
// 						&element.Message, &element.IsEdited, &element.Forum, &element.Thread, 
// 						&element.Created); err != nil {
// 						return nil, err
// 					}
// 					posts = append(posts, &element)
// 				}
// 			}
// 		} else {
// 			if since != 0 {
// 				query := `SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts
// 				  WHERE thread = $1 and post_tree > (select post_tree from posts where id = $2)
// 				  ORDER BY post_tree, id
// 				  LIMIT $3;`

// 				rows, err := dbPost.pool.Query(context.Background(), query, id, since, limit)
// 				if err != nil {
// 					return nil, err
// 				}

// 				for rows.Next() {
// 					element := models.Post{}
// 					if err := rows.Scan(&element.Id, &element.Parent, &element.Author, 
// 						&element.Message, &element.IsEdited, &element.Forum, &element.Thread, 
// 						&element.Created); err != nil {
// 						return nil, err
// 					}
// 					posts = append(posts, &element)
// 				}
// 			} else {
// 				query := `SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts
// 				  WHERE thread = $1
// 				  ORDER BY post_tree
// 				  LIMIT $2;`

// 				rows, err := dbPost.pool.Query(context.Background(), query, id, limit)
// 				if err != nil {
// 					return nil, err
// 				}

// 				for rows.Next() {
// 					element := models.Post{}
// 					if err := rows.Scan(&element.Id, &element.Parent, &element.Author, 
// 						&element.Message, &element.IsEdited, &element.Forum, &element.Thread, 
// 						&element.Created); err != nil {
// 						return nil, err
// 					}
// 					posts = append(posts, &element)
// 				}
// 			}
// 		}
// 	} else if sort == "parent_tree" {
// 		if desc {
// 			if since != 0 {
// 				query := `SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts
// 				  WHERE post_tree[1] in (
// 				  	select id from posts where thread = $1 and parent = 0 and id < (select post_tree[1] from posts where id = $2) order by id DESC limit $3
// 				  )
// 				  ORDER BY post_tree[1] DESC, post_tree, id`

// 				rows, err := dbPost.pool.Query(context.Background(), query, id, since, limit)
// 				if err != nil {
// 					return nil, err
// 				}

// 				for rows.Next() {
// 					element := models.Post{}
// 					if err := rows.Scan(&element.Id, &element.Parent, &element.Author, 
// 						&element.Message, &element.IsEdited, &element.Forum, &element.Thread, 
// 						&element.Created); err != nil {
// 						return nil, err
// 					}
// 					posts = append(posts, &element)
// 				}
// 			} else {
// 				query := `SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts
// 				  WHERE post_tree[1] in (
// 				  	select id from posts where parent = 0 and thread = $1 order by id DESC limit $2
// 				  )
// 				  ORDER BY post_tree[1] DESC, post_tree`

// 				rows, err := dbPost.pool.Query(context.Background(), query, id, limit)
// 				if err != nil {
// 					return nil, err
// 				}

// 				for rows.Next() {
// 					element := models.Post{}
// 					if err := rows.Scan(&element.Id, &element.Parent, &element.Author, 
// 						&element.Message, &element.IsEdited, &element.Forum, &element.Thread, 
// 						&element.Created); err != nil {
// 						return nil, err
// 					}
// 					posts = append(posts, &element)
// 				}
// 			}
// 		} else {
// 			if since != 0 {
// 				query := `SELECT id, parent, author, message, is_edited, forum, thread, created FROM posts
// 				  WHERE post_tree[1] in (
// 				  	select id from posts where thread = $1 and parent = 0 and id > (select post_tree[1] from posts where id = $2) order by id limit $3
// 				  )
// 				  ORDER BY post_tree, id`

// 				rows, err := dbPost.pool.Query(context.Background(), query, id, since, limit)
// 				if err != nil {
// 					return nil, err
// 				}

// 				for rows.Next() {
// 					element := models.Post{}
// 					if err := rows.Scan(&element.Id, &element.Parent, &element.Author, 
// 						&element.Message, &element.IsEdited, &element.Forum, &element.Thread, 
// 						&element.Created); err != nil {
// 						return nil, err
// 					}
// 					posts = append(posts, &element)
// 				}
// 			} else {
// 				query := `SELECT id, parent, author, message, is_edited, forum, thread, created
// 				  FROM posts
// 				  WHERE post_tree[1] in (
// 				  	select id from posts where parent = 0 and thread = $1 order by id limit $2
// 				  )
// 				  ORDER BY post_tree, id`

// 				rows, err := dbPost.pool.Query(context.Background(), query, id, limit)
// 				if err != nil {
// 					return nil, err
// 				}

// 				for rows.Next() {
// 					element := models.Post{}
// 					if err := rows.Scan(&element.Id, &element.Parent, &element.Author, 
// 						&element.Message, &element.IsEdited, &element.Forum, &element.Thread, 
// 						&element.Created); err != nil {
// 						return nil, err
// 					}
// 					posts = append(posts, &element)
// 				}
// 			}
// 		}
// 	}

// 	return posts, nil
// }

