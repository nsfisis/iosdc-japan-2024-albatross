// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Game struct {
	GameID    int32
	Type      string
	CreatedAt pgtype.Timestamp
	State     string
}

type User struct {
	UserID          int32
	Username        string
	DisplayUsername string
	IconUrl         pgtype.Text
	IsAdmin         bool
	CreatedAt       pgtype.Timestamp
}

type UserAuth struct {
	UserAuthID   int32
	UserID       int32
	AuthType     string
	PasswordHash pgtype.Text
}
