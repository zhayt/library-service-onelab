package model

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

type contextKey string

const ContextUserID = contextKey("userID")
