package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/nsfisis/iosdc-2024-albatross-backend/db"
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
	var iconPath *string
	if user.IconPath.Valid {
		iconPath = &user.IconPath.String
	}
	claims := &JWTClaims{
		UserID:      int(user.UserID),
		Username:    user.Username,
		DisplayName: user.DisplayName,
		IconPath:    iconPath,
		IsAdmin:     user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
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
