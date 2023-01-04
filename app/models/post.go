package models

import "time"

type Post struct {
	Id       uint64    `json:"id,omitempty" gorm:"column:id"`
	Parent   uint64    `json:"parent,omitempty" gorm:"column:parent"`
	Author   string    `json:"author,omitempty" gorm:"column:author"`
	Message  string    `json:"message,omitempty" gorm:"column:message"`
	IsEdited bool      `json:"isEdited,omitempty" gorm:"column:is_edited"`
	Forum    string    `json:"forum,omitempty" gorm:"column:forum"`
	Thread   uint64    `json:"thread,omitempty" gorm:"column:thread"`
	Created  time.Time `json:"created,omitempty" gorm:"column:created"`
}

type PostDetails struct {
	Post   *Post
	User   *User
	Forum  *Forum
	Thread *Thread
}
