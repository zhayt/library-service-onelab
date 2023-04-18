package middleware

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/zhayt/user-storage-service/config"
	"github.com/zhayt/user-storage-service/internal/model"
	"net/http"
	"strings"
	"time"
)

type JWTAuth struct {
	jwtKey []byte
}

func NewJWTAuth(cfg *config.Config) *JWTAuth {
	return &JWTAuth{jwtKey: []byte(cfg.JWTKey)}
}

func (m *JWTAuth) GenerateJWT(username string, userID int) (tokenString string, err error) {
	expirationTime := time.Now().Add(12 * time.Hour)

	claims := &model.JWTClaim{
		Username: username,
		UserID:   userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.jwtKey)
}

func (m *JWTAuth) ValidateToken(accessToken string) (*model.JWTClaim, error) {
	token, err := jwt.ParseWithClaims(accessToken, &model.JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid sign-in method")
			}

			return m.jwtKey, nil
		})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*model.JWTClaim)
	if !ok {
		return nil, errors.New("couldn't parse claims")
	}
	if time.Now().Local().After(time.Unix(claims.ExpiresAt, 0)) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

func (m *JWTAuth) ValidateAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		token := extractToken(e.Request())
		if token != "" {
			claims, err := m.ValidateToken(token)
			if err != nil {
				return echo.NewHTTPError(403, err.Error())
			}

			ctx := context.WithValue(e.Request().Context(), model.ContextUserID, claims.UserID)
			ctx = context.WithValue(ctx, model.ContextUserName, claims.Username)
			e.SetRequest(e.Request().WithContext(ctx))
		}

		return next(e)
	}
}

func extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}

	return ""
}
