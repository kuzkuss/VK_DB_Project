package threads

import (
	"errors"

	forums "github.com/kuzkuss/VK_DB_Project/internal/app/models/forumsRepository"
	users "github.com/kuzkuss/VK_DB_Project/internal/app/models/usersRepository"
)

type ThreadsStore interface {
	CreateThread(t Thread) (*Thread, error)
	SelectThread(slug string) (*Thread, error)
}

type ThreadsRep struct {
	threadsStore ThreadsStore
	usersStore   users.UsersStore
	forumsStore  forums.ForumsStore
}

type Thread struct {
	Id      int    `json: id`
	Title   string `json: title`
	Author  string `json: author`
	Forum   string `json: forum`
	Message string `json: message`
	Votes   int    `json: votes`
	Slug    string `json: slug`
	Created string `json: created`
}

func NewThreadsRep(ts ThreadsStore, us users.UsersStore, fs forums.ForumsStore) *ThreadsRep {
	return &ThreadsRep{
		threadsStore: ts,
		usersStore:   us,
		forumsStore:  fs,
	}
}

func (tr *ThreadsRep) CreateThread(t Thread, slug string) (*Thread, error) {
	if _, err := tr.usersStore.SelectUser(t.Author); err != nil {
		return nil, errors.New("Can't find user with name " + t.Author)
	}

	if _, err := tr.forumsStore.SelectForum(slug); err != nil {
		return nil, errors.New("Can't find forum with slug " + slug)
	}

	if trd, _ := tr.threadsStore.SelectThread(t.Slug); trd != nil {
		return trd, errors.New("Thread existed")
	}

	nf, err := tr.threadsStore.CreateThread(t)

	if err != nil {
		return nil, errors.New("Create thread error")
	}

	return nf, nil
}

// func (fr *ForumsRep) GetForum(slug string) (*Forum, error) {
// 	frm, err := fr.forumsStore.SelectForum(slug)
// 	if err != nil {
// 		return nil, errors.New("Can't find forum with slug " + slug)
// 	}

// 	return frm, nil
// }
