package admin

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"

	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/auth"
	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/db"
)

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

type AdminHandler struct {
	q    *db.Queries
	hubs GameHubsInterface
}

type GameHubsInterface interface {
	StartGame(gameID int) error
}

func NewAdminHandler(q *db.Queries, hubs GameHubsInterface) *AdminHandler {
	return &AdminHandler{
		q:    q,
		hubs: hubs,
	}
}

func newAdminMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			jwt, err := c.Cookie("albatross_token")
			if err != nil {
				return c.Redirect(http.StatusSeeOther, "/login")
			}
			claims, err := auth.ParseJWT(jwt.Value)
			if err != nil {
				return c.Redirect(http.StatusSeeOther, "/login")
			}
			if !claims.IsAdmin {
				return echo.NewHTTPError(http.StatusForbidden)
			}
			return next(c)
		}
	}
}

func (h *AdminHandler) RegisterHandlers(g *echo.Group) {
	g.Use(newAssetsMiddleware())
	g.Use(newAdminMiddleware())

	g.GET("/dashboard", h.getDashboard)
	g.GET("/users", h.getUsers)
	g.GET("/users/:userID", h.getUserEdit)
	g.GET("/games", h.getGames)
	g.GET("/games/:gameID", h.getGameEdit)
	g.POST("/games/:gameID", h.postGameEdit)
}

func (h *AdminHandler) getDashboard(c echo.Context) error {
	return c.Render(http.StatusOK, "dashboard", echo.Map{
		"Title": "Dashboard",
	})
}

func (h *AdminHandler) getUsers(c echo.Context) error {
	rows, err := h.q.ListUsers(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	users := make([]echo.Map, len(rows))
	for i, u := range rows {
		users[i] = echo.Map{
			"UserID":      u.UserID,
			"Username":    u.Username,
			"DisplayName": u.DisplayName,
			"IconPath":    u.IconPath,
			"IsAdmin":     u.IsAdmin,
		}
	}

	return c.Render(http.StatusOK, "users", echo.Map{
		"Title": "Users",
		"Users": users,
	})
}

func (h *AdminHandler) getUserEdit(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user id")
	}
	row, err := h.q.GetUserByID(c.Request().Context(), int32(userID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.Render(http.StatusOK, "user_edit", echo.Map{
		"Title": "User Edit",
		"User": echo.Map{
			"UserID":      row.UserID,
			"Username":    row.Username,
			"DisplayName": row.DisplayName,
			"IconPath":    row.IconPath,
			"IsAdmin":     row.IsAdmin,
		},
	})
}

func (h *AdminHandler) getGames(c echo.Context) error {
	rows, err := h.q.ListGames(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	games := make([]echo.Map, len(rows))
	for i, g := range rows {
		var startedAt string
		if !g.StartedAt.Valid {
			startedAt = g.StartedAt.Time.In(jst).Format("2006-01-02T15:04")
		}
		games[i] = echo.Map{
			"GameID":          g.GameID,
			"GameType":        g.GameType,
			"State":           g.State,
			"DisplayName":     g.DisplayName,
			"DurationSeconds": g.DurationSeconds,
			"StartedAt":       startedAt,
			"ProblemID":       g.ProblemID,
		}
	}

	return c.Render(http.StatusOK, "games", echo.Map{
		"Title": "Games",
		"Games": games,
	})
}

func (h *AdminHandler) getGameEdit(c echo.Context) error {
	gameID, err := strconv.Atoi(c.Param("gameID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid game id")
	}
	row, err := h.q.GetGameByID(c.Request().Context(), int32(gameID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	var startedAt string
	if !row.StartedAt.Valid {
		startedAt = row.StartedAt.Time.In(jst).Format("2006-01-02T15:04")
	}

	return c.Render(http.StatusOK, "game_edit", echo.Map{
		"Title": "Game Edit",
		"Game": echo.Map{
			"GameID":          row.GameID,
			"GameType":        row.GameType,
			"State":           row.State,
			"DisplayName":     row.DisplayName,
			"DurationSeconds": row.DurationSeconds,
			"StartedAt":       startedAt,
			"ProblemID":       row.ProblemID,
		},
	})
}

func (h *AdminHandler) postGameEdit(c echo.Context) error {
	gameID, err := strconv.Atoi(c.Param("gameID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid game id")
	}
	row, err := h.q.GetGameByID(c.Request().Context(), int32(gameID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	gameType := c.FormValue("game_type")
	state := c.FormValue("state")
	displayName := c.FormValue("display_name")
	durationSeconds, err := strconv.Atoi(c.FormValue("duration_seconds"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid duration_seconds")
	}
	var problemID *int
	{
		problemIDRaw := c.FormValue("problem_id")
		if problemIDRaw != "" {
			problemIDInt, err := strconv.Atoi(problemIDRaw)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid problem_id")
			}
			problemID = &problemIDInt
		}
	}
	var startedAt *time.Time
	{
		startedAtRaw := c.FormValue("started_at")
		if startedAtRaw != "" {
			startedAtTime, err := time.Parse("2006-01-02T15:04", startedAtRaw)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid started_at")
			}
			startedAt = &startedAtTime
		}
	}

	var changedStartedAt pgtype.Timestamp
	if startedAt == nil {
		changedStartedAt = pgtype.Timestamp{
			Valid: false,
		}
	} else {
		changedStartedAt = pgtype.Timestamp{
			Time:  *startedAt,
			Valid: true,
		}
	}
	var changedProblemID *int32
	if problemID == nil {
		changedProblemID = nil
	} else {
		changedProblemID = new(int32)
		*changedProblemID = int32(*problemID)
	}

	{
		// TODO:
		if state != row.State && state == "prepare" {
			h.hubs.StartGame(int(gameID))
		}
	}

	err = h.q.UpdateGame(c.Request().Context(), db.UpdateGameParams{
		GameID:          int32(gameID),
		GameType:        gameType,
		State:           state,
		DisplayName:     displayName,
		DurationSeconds: int32(durationSeconds),
		StartedAt:       changedStartedAt,
		ProblemID:       changedProblemID,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
