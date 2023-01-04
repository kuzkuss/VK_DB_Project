package repository

import (
	"github.com/kuzkuss/VK_DB_Project/app/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RepositoryI interface {
	CreateForum(forum *models.Forum) (error)
	SelectForumBySlug(slug string) (*models.Forum, error)
	SelectForumUsers(slug string, limit int, since string, desc bool) ([]*models.User, error)
	CreateForumUser(forum string, user string) (error)
}

type dataBase struct {
	db *gorm.DB
}

func New(db *gorm.DB) RepositoryI {
	return &dataBase{
		db: db,
	}
}

func (dbForum *dataBase) CreateForum(forum *models.Forum) (error) {
	tx := dbForum.db.Omit("posts", "threads").Create(forum)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table forums)")
	}

	return nil
}

func (dbForum *dataBase) SelectForumBySlug(slug string) (*models.Forum, error) {
	forum := models.Forum{}

	tx := dbForum.db.Where("slug = ?", slug).Take(&forum)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table forums)")
	}

	return &forum, nil
}

func (dbForum *dataBase) SelectForumUsers(slug string, limit int, since string, desc bool) ([]*models.User, error) {
	users := make([]*models.User, 0, 10)

	if desc {
		if since != "" {
			tx := dbForum.db.Table("forum_user").Select("u.nickname, u.fullname, u.about, u.email").Limit(limit).Where("forum = ? AND user < ?", slug, since).Joins("JOIN user u ON u.nickname=forum_user.user").Order("nickname desc").Scan(&users)
			if tx.Error != nil {
				return nil, errors.Wrap(tx.Error, "database error (table forum_user)")
			}
		} else {
			tx := dbForum.db.Table("forum_user").Select("u.nickname, u.fullname, u.about, u.email").Limit(limit).Where("forum = ?", slug).Joins("JOIN user u ON u.nickname=forum_user.user").Order("nickname desc").Scan(&users)
			if tx.Error != nil {
				return nil, errors.Wrap(tx.Error, "database error (table forum_user)")
			}
		}
	} else {
		if since != "" {
			tx := dbForum.db.Table("forum_user").Select("u.nickname, u.fullname, u.about, u.email").Limit(limit).Where("forum = ? AND user > ?", slug, since).Joins("JOIN user u ON u.nickname=forum_user.user").Order("nickname").Scan(&users)
			if tx.Error != nil {
				return nil, errors.Wrap(tx.Error, "database error (table forum_user)")
			}
		} else {
			tx := dbForum.db.Table("forum_user").Select("u.nickname, u.fullname, u.about, u.email").Limit(limit).Where("forum = ?", slug).Joins("JOIN user u ON u.nickname=forum_user.user").Order("nickname").Scan(&users)
			if tx.Error != nil {
				return nil, errors.Wrap(tx.Error, "database error (table forum_user)")
			}
		}
	}

	return users, nil
}

func (dbForum *dataBase) CreateForumUser(forum string, user string) (error) {
	fu := models.ForumUser {
		Forum: forum,
		User: user,
	}
	tx := dbForum.db.Table("forum_user").Clauses(clause.OnConflict{DoNothing: true}).Create(&fu)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table forum_user)")
	}

	return nil
}

