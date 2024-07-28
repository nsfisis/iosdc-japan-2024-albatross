package game

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type sockHandler struct {
	hubs *GameHubs
}

func newSockHandler(hubs *GameHubs) *sockHandler {
	return &sockHandler{
		hubs: hubs,
	}
}

func (h *sockHandler) HandleSockGolfPlay(c echo.Context) error {
	// TODO: auth
	gameId := c.Param("gameId")
	gameIdInt, err := strconv.Atoi(gameId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid game id")
	}
	var foundHub *gameHub
	for _, hub := range h.hubs.hubs {
		if hub.game.gameID == gameIdInt {
			foundHub = hub
			break
		}
	}
	if foundHub == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Game not found")
	}
	return servePlayerWs(foundHub, c.Response(), c.Request(), "a")
}

func (h *sockHandler) HandleSockGolfWatch(c echo.Context) error {
	// TODO: auth
	gameId := c.Param("gameId")
	gameIdInt, err := strconv.Atoi(gameId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid game id")
	}
	var foundHub *gameHub
	for _, hub := range h.hubs.hubs {
		if hub.game.gameID == gameIdInt {
			foundHub = hub
			break
		}
	}
	if foundHub == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Game not found")
	}
	return serveWatcherWs(foundHub, c.Response(), c.Request())
}
