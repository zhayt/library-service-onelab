package model

import "github.com/dgrijalva/jwt-go"

type JWTClaim struct {
	Username string `json:"username"`
	UserID   int    `json:"userID"`
	jwt.StandardClaims
}
