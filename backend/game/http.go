package game

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/nsfisis/iosdc-2024-albatross/backend/auth"
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
	jwt := c.QueryParam("token")
	claims, err := auth.ParseJWT(jwt)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	// TODO: check user permission

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
	return servePlayerWs(foundHub, c.Response(), c.Request(), claims.UserID)
}

func (h *sockHandler) HandleSockGolfWatch(c echo.Context) error {
	jwt := c.QueryParam("token")
	claims, err := auth.ParseJWT(jwt)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	if !claims.IsAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "Permission denied")
	}

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
