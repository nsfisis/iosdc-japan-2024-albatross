package auth

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/db"
)

func Login(ctx context.Context, queries *db.Queries, username, password string) (int, error) {
	userAuth, err := queries.GetUserAuthByUsername(ctx, username)
	if err != nil {
		return 0, err
	}
	if userAuth.AuthType == "password" {
		passwordHash := userAuth.PasswordHash
		if passwordHash == nil {
			panic("inconsistant data")
		}
		err := bcrypt.CompareHashAndPassword([]byte(*passwordHash), []byte(password))
		if err != nil {
			return 0, err
		}
		return int(userAuth.UserID), nil
	}
	return 0, fmt.Errorf("not implemented")
}
