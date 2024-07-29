package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/nsfisis/iosdc-2024-albatross/backend/db"
)

type JWTClaims struct {
	UserID      int     `json:"user_id"`
	Username    string  `json:"username"`
	DisplayName string  `json:"display_name"`
	IconPath    *string `json:"icon_path"`
	IsAdmin     bool    `json:"is_admin"`
	jwt.RegisteredClaims
}

func NewJWT(user *db.User) (string, error) {
	claims := &JWTClaims{
		UserID:      int(user.UserID),
		Username:    user.Username,
		DisplayName: user.DisplayName,
		IconPath:    user.IconPath,
		IsAdmin:     user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("TODO"))
}

func NewShortLivedJWT(claims *JWTClaims) (string, error) {
	newClaims := &JWTClaims{
		UserID:      claims.UserID,
		Username:    claims.Username,
		DisplayName: claims.DisplayName,
		IconPath:    claims.IconPath,
		IsAdmin:     claims.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	return token.SignedString([]byte("TODO"))
}

func ParseJWT(token string) (*JWTClaims, error) {
	claims := new(JWTClaims)
	t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("TODO"), nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
