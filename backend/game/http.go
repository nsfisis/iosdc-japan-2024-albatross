package game

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/auth"
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

	gameID, err := strconv.Atoi(c.Param("gameId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid game id")
	}
	hub := h.hubs.getHub(gameID)
	if hub == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Game not found")
	}
	return servePlayerWs(hub, c.Response(), c.Request(), claims.UserID)
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

	gameID, err := strconv.Atoi(c.Param("gameId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid game id")
	}
	hub := h.hubs.getHub(gameID)
	if hub == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Game not found")
	}
	return serveWatcherWs(hub, c.Response(), c.Request())
}
