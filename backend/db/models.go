// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Game struct {
	GameID          int32
	GameType        string
	State           string
	DisplayName     string
	DurationSeconds int32
	CreatedAt       pgtype.Timestamp
	StartedAt       pgtype.Timestamp
	ProblemID       *int32
}

type GamePlayer struct {
	GameID int32
	UserID int32
}

type Problem struct {
	ProblemID   int32
	Title       string
	Description string
}

type Submission struct {
	SubmissionID int32
	GameID       int32
	UserID       int32
	Code         string
	CodeSize     int32
	CreatedAt    pgtype.Timestamp
}

type Testcase struct {
	TestcaseID int32
	ProblemID  int32
	Stdin      string
	Stdout     string
}

type TestcaseResult struct {
	TestcaseResultID int32
	SubmissionID     int32
	TestcaseID       int32
	Status           string
	Stdout           string
	Stderr           string
}

type User struct {
	UserID      int32
	Username    string
	DisplayName string
	IconPath    *string
	IsAdmin     bool
	CreatedAt   pgtype.Timestamp
}

type UserAuth struct {
	UserAuthID   int32
	UserID       int32
	AuthType     string
	PasswordHash *string
}
