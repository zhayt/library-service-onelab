package model

type User struct {
	FIO          string `json:"fio"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
}
