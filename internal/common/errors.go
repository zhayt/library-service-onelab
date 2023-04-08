package common

import "errors"
// зачем этот пакет ?
var (
	ErrInvalidData   = errors.New("invalid data")
	ErrNameTaken     = errors.New("name taken")
	ErrUserNotExists = errors.New("user not exists")
	ErrEmptyField    = errors.New("empty field")
)
