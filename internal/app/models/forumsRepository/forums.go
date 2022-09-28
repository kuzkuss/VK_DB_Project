package forums

import (
	"errors"

	users "github.com/kuzkuss/VK_DB_Project/internal/app/models/usersRepository"
)

type ForumsStore interface {
	CreateForum(f Forum) (*Forum, error)
	SelectForum(slug string) (*Forum, error)
}

type ForumsRep struct {
	forumsStore ForumsStore
	usersStore users.UsersStore
}

type Forum struct {
	Title   string `json: title`
	User    string `json: user`
	Slug    string `json: slug`
	Posts   int    `json: posts`
	Threads int    `json: threads`
}

func NewForumsRep(fs ForumsStore, us users.UsersStore) *ForumsRep {
	return &ForumsRep{
		forumsStore: fs,
		usersStore: us,
	}
}

func (fr *ForumsRep) CreateForum(f Forum) (*Forum, error) {
	if _, err := fr.usersStore.SelectUser(f.User); err != nil {
		return nil, errors.New("Can't find user with name " + f.User)
	}

	if frm, _ := fr.forumsStore.SelectForum(f.Slug); frm != nil {
		return frm, errors.New("Forum existed")
	}

	nf, err := fr.forumsStore.CreateForum(f)

	if err != nil {
		return nil, errors.New("Create forum error")
	}

	return nf, nil
}

func (fr *ForumsRep) GetForum(slug string) (*Forum, error) {
	frm, err := fr.forumsStore.SelectForum(slug)
	if err != nil {
		return nil, errors.New("Can't find forum with slug " + slug)
	}

	return frm, nil
}

