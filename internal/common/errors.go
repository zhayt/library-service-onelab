package common

import "errors"

var (
	ErrInvalidData   = errors.New("invalid data")
	ErrNameTaken     = errors.New("name taken")
	ErrUserNotExists = errors.New("user not exists")
	ErrNoRows        = errors.New("empty table")
	ErrEmptyField    = errors.New("empty field")
)

var (
	ErrBookNotExists = errors.New("book not exists")
)
