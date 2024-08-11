package game

import (
	"time"

	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/api"
)

type gameType = api.GameGameType
type gameState = api.GameState

const (
	gameType1v1         = api.N1V1
	gameTypeMultiplayer = api.Multiplayer

	gameStateClosed   gameState = api.Closed
	gameStateWaiting  gameState = api.Waiting
	gameStateStarting gameState = api.Starting
	gameStateGaming   gameState = api.Gaming
	gameStateFinished gameState = api.Finished
)

type game struct {
	gameID          int
	gameType        gameType
	state           gameState
	displayName     string
	durationSeconds int
	startedAt       *time.Time
	problem         *problem
	playerCount     int
}

type problem struct {
	problemID   int
	title       string
	description string
}
