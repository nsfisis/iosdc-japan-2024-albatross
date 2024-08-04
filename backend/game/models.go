package game

import (
	"time"

	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/api"
)

type gameState = api.GameState

const (
	gameStateClosed         gameState = api.Closed
	gameStateWaitingEntries gameState = api.WaitingEntries
	gameStateWaitingStart   gameState = api.WaitingStart
	gameStatePrepare        gameState = api.Prepare
	gameStateStarting       gameState = api.Starting
	gameStateGaming         gameState = api.Gaming
	gameStateFinished       gameState = api.Finished
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
