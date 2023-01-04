package models

import (
	"github.com/pkg/errors"
)

var (
	ErrNotFound            = errors.New("item is not found")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrConflict            = errors.New("item already exists")
	ErrBadRequest          = errors.New("bad request")
	ErrConflictFriend      = errors.New("friend already exists")
	ErrInternalServerError = errors.New("internal server error")
	ErrPermissionDenied    = errors.New("permission denied")
)
