package model

type User struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Fatherland   string `json:"fatherland"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
}
