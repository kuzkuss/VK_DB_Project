package repository

import (
	"github.com/kuzkuss/VK_DB_Project/app/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type RepositoryI interface {
	ClearData() (error)
	SelectStatus() (*models.ServiceStatus, error)
}

type dataBase struct {
	db *gorm.DB
}

func New(db *gorm.DB) RepositoryI {
	return &dataBase{
		db: db,
	}
}

func (dbService *dataBase) ClearData() (error) {
	tx := dbService.db.Exec("TRUNCATE post, thread, forum, user, forum_user cascade;")
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error")
	}

	return nil
}

func (dbService *dataBase) SelectStatus() (*models.ServiceStatus, error) {
	status := models.ServiceStatus{}

	var count int64
	tx := dbService.db.Model(&models.User{}).Count(&count)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table user)")
	}
	status.UserCount = count
	tx = dbService.db.Model(&models.Forum{}).Count(&count)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table forum)")
	}
	status.ForumCount = count
	tx = dbService.db.Model(&models.Thread{}).Count(&count)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table post)")
	}
	status.ThreadCount = count
	tx = dbService.db.Model(&models.Post{}).Count(&count)
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table thread)")
	}
	status.PostCount = count

	return &status, nil
}

