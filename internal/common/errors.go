package common

import "errors"

var (
	ErrInvalidData   = errors.New("invalid data")
	ErrNameTaken     = errors.New("name taken")
	ErrUserNotExists = errors.New("user not exists")
	ErrEmptyField    = errors.New("empty field")
)
