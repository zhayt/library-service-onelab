package model

import "time"

type User struct {
	ID       int    `json:"id"`
	FIO      string `json:"fio"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdateFIO struct {
	ID  int
	FIO string `json:"fio"`
}

type UserUpdatePassword struct {
	ID                int
	CurrentPassword   string `json:"currentPassword"`
	NewPassword       string `json:"newPassword"`
	NewPasswordRepeat string `json:"newPasswordRepeat"`
}

type UserAccount struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	CardNumber string    `json:"cardNumber"`
	Balance    float64   `json:"currentBalance"`
	CreatedAt  time.Time `json:"createdAt"`
}

type NewUserAccountBalance struct {
	AccountID           uint    `json:"accountID"`
	Name                string  `json:"name"`
	CardNumber          string  `json:"cardNumber"`
	ReplenishmentAmount float64 `json:"replenishmentAmount"`
}

type contextKey string

const ContextUserID = contextKey("userID")
const ContextUserName = contextKey("userName")
