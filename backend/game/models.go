package game

import (
	"time"

	"github.com/nsfisis/iosdc-2024-albatross/backend/api"
)

type gameState = api.GameState

const (
	gameStateClosed         gameState = api.GameStateClosed
	gameStateWaitingEntries gameState = api.GameStateWaitingEntries
	gameStateWaitingStart   gameState = api.GameStateWaitingStart
	gameStatePrepare        gameState = api.GameStatePrepare
	gameStateStarting       gameState = api.GameStateStarting
	gameStateGaming         gameState = api.GameStateGaming
	gameStateFinished       gameState = api.GameStateFinished
)

type game struct {
	gameID          int
	state           gameState
	displayName     string
	durationSeconds int
	startedAt       *time.Time
	problem         *problem
}

type problem struct {
	problemID   int
	title       string
	description string
}
