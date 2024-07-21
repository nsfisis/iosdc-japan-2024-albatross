package main

import (
	"context"
	"fmt"

	"github.com/nsfisis/iosdc-2024-albatross-backend/db"
)

func authLogin(ctx context.Context, queries *db.Queries, username, password string) (int, error) {
	userAuth, err := queries.GetUserAuthFromUsername(ctx, username)
	if err != nil {
		return 0, err
	}
	if userAuth.AuthType == "bypass" {
		return int(userAuth.UserID), nil
	}
	return 0, fmt.Errorf("not implemented")
}
