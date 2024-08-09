package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/db"
)

var (
	jwtSecret []byte
)

func init() {
	jwtSecret = []byte(os.Getenv("ALBATROSS_JWT_SECRET"))
	if len(jwtSecret) == 0 {
		panic("ALBATROSS_JWT_SECRET is not set")
	}
}

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
	return token.SignedString(jwtSecret)
}

func NewAnonymousJWT() (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
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
	return token.SignedString(jwtSecret)
}

func ParseJWT(token string) (*JWTClaims, error) {
	claims := new(JWTClaims)
	t, err := jwt.ParseWithClaims(token, claims, func(*jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
