package users

type User struct {
	Nickname   string
	Fullname    string
	About    string
	Email   string
}

type UsersStore interface {
	SelectUser(name string) (*User, error)
}

type UsersRep struct {
	usersStore UsersStore
}

func (*UsersRep) SelectUser(name string) (*User, error) {
	return nil, nil
}

