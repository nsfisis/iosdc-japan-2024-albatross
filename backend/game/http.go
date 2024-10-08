package game

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/auth"
)

type SockHandler struct {
	hubs *Hubs
}

func newSockHandler(hubs *Hubs) *SockHandler {
	return &SockHandler{
		hubs: hubs,
	}
}

func (h *SockHandler) HandleSockGolfPlay(c echo.Context) error {
	jwt := c.QueryParam("token")
	claims, err := auth.ParseJWT(jwt)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	// TODO: check user permission

	gameID, err := strconv.Atoi(c.Param("gameID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid game id")
	}
	hub := h.hubs.getHub(gameID)
	if hub == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Game not found")
	}
	return servePlayerWs(hub, c.Response(), c.Request(), claims.UserID)
}

func (h *SockHandler) HandleSockGolfWatch(c echo.Context) error {
	jwt := c.QueryParam("token")
	claims, err := auth.ParseJWT(jwt)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	if !claims.IsAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "Permission denied")
	}

	gameID, err := strconv.Atoi(c.Param("gameID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid game id")
	}
	hub := h.hubs.getHub(gameID)
	if hub == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Game not found")
	}

	if hub.game.gameType != gameType1v1 {
		return echo.NewHTTPError(http.StatusBadRequest, "Only 1v1 game is supported")
	}

	return serveWatcherWs(hub, c.Response(), c.Request())
}
