package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

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

func NewJWTMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JWTClaims)
		},
		SigningKey: []byte("TODO"),
	})
}

func GetJWTClaimsFromEchoContext(c echo.Context) *JWTClaims {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)
	return claims
}
