package api

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"

	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/auth"
	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/db"
)

var _ StrictServerInterface = (*ApiHandler)(nil)

type ApiHandler struct {
	q *db.Queries
}

func NewHandler(queries *db.Queries) *ApiHandler {
	return &ApiHandler{
		q: queries,
	}
}

func (h *ApiHandler) GetAdminGames(ctx context.Context, request GetAdminGamesRequestObject) (GetAdminGamesResponseObject, error) {
	user := ctx.Value("user").(*auth.JWTClaims)
	if !user.IsAdmin {
		return GetAdminGames403JSONResponse{
			Message: "Forbidden",
		}, nil
	}
	gameRows, err := h.q.ListGames(ctx)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	games := make([]Game, len(gameRows))
	for i, row := range gameRows {
		var startedAt *int
		if row.StartedAt.Valid {
			startedAtTimestamp := int(row.StartedAt.Time.Unix())
			startedAt = &startedAtTimestamp
		}
		var problem *Problem
		if row.ProblemID != nil {
			if row.Title == nil || row.Description == nil {
				panic("inconsistent data")
			}
			problem = &Problem{
				ProblemId:   int(*row.ProblemID),
				Title:       *row.Title,
				Description: *row.Description,
			}
		}
		games[i] = Game{
			GameId:          int(row.GameID),
			State:           GameState(row.State),
			DisplayName:     row.DisplayName,
			DurationSeconds: int(row.DurationSeconds),
			StartedAt:       startedAt,
			Problem:         problem,
		}
	}
	return GetAdminGames200JSONResponse{
		Games: games,
	}, nil
}

func (h *ApiHandler) GetAdminGamesGameId(ctx context.Context, request GetAdminGamesGameIdRequestObject) (GetAdminGamesGameIdResponseObject, error) {
	user := ctx.Value("user").(*auth.JWTClaims)
	if !user.IsAdmin {
		return GetAdminGamesGameId403JSONResponse{
			Message: "Forbidden",
		}, nil
	}
	gameId := request.GameId
	row, err := h.q.GetGameById(ctx, int32(gameId))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return GetAdminGamesGameId404JSONResponse{
				Message: "Game not found",
			}, nil
		} else {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	var startedAt *int
	if row.StartedAt.Valid {
		startedAtTimestamp := int(row.StartedAt.Time.Unix())
		startedAt = &startedAtTimestamp
	}
	var problem *Problem
	if row.ProblemID != nil {
		if row.Title == nil || row.Description == nil {
			panic("inconsistent data")
		}
		problem = &Problem{
			ProblemId:   int(*row.ProblemID),
			Title:       *row.Title,
			Description: *row.Description,
		}
	}
	game := Game{
		GameId:          int(row.GameID),
		State:           GameState(row.State),
		DisplayName:     row.DisplayName,
		DurationSeconds: int(row.DurationSeconds),
		StartedAt:       startedAt,
		Problem:         problem,
	}
	return GetAdminGamesGameId200JSONResponse{
		Game: game,
	}, nil
}

func (h *ApiHandler) PutAdminGamesGameId(ctx context.Context, request PutAdminGamesGameIdRequestObject) (PutAdminGamesGameIdResponseObject, error) {
	user := ctx.Value("user").(*auth.JWTClaims)
	if !user.IsAdmin {
		return PutAdminGamesGameId403JSONResponse{
			Message: "Forbidden",
		}, nil
	}
	gameID := request.GameId
	displayName := request.Body.DisplayName
	durationSeconds := request.Body.DurationSeconds
	problemID := request.Body.ProblemId
	startedAt := request.Body.StartedAt
	state := request.Body.State

	game, err := h.q.GetGameById(ctx, int32(gameID))
	if err != nil {
		if err == pgx.ErrNoRows {
			return PutAdminGamesGameId404JSONResponse{
				Message: "Game not found",
			}, nil
		} else {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	var changedState string
	if state != nil {
		changedState = string(*state)
	} else {
		changedState = game.State
	}
	var changedDisplayName string
	if displayName != nil {
		changedDisplayName = *displayName
	} else {
		changedDisplayName = game.DisplayName
	}
	var changedDurationSeconds int32
	if durationSeconds != nil {
		changedDurationSeconds = int32(*durationSeconds)
	} else {
		changedDurationSeconds = game.DurationSeconds
	}
	var changedStartedAt pgtype.Timestamp
	if startedAt != nil {
		startedAtValue, err := startedAt.Get()
		if err == nil {
			changedStartedAt = pgtype.Timestamp{
				Time:  time.Unix(int64(startedAtValue), 0),
				Valid: true,
			}
		}
	} else {
		changedStartedAt = game.StartedAt
	}
	var changedProblemID *int32
	if problemID != nil {
		problemIDValue, err := problemID.Get()
		if err == nil {
			changedProblemID = new(int32)
			*changedProblemID = int32(problemIDValue)
		}
	} else {
		changedProblemID = game.ProblemID
	}

	err = h.q.UpdateGame(ctx, db.UpdateGameParams{
		GameID:          int32(gameID),
		State:           changedState,
		DisplayName:     changedDisplayName,
		DurationSeconds: changedDurationSeconds,
		StartedAt:       changedStartedAt,
		ProblemID:       changedProblemID,
	})
	if err != nil {
		return PutAdminGamesGameId400JSONResponse{
			Message: err.Error(),
		}, nil
	}

	return PutAdminGamesGameId204Response{}, nil
}

func (h *ApiHandler) GetAdminUsers(ctx context.Context, request GetAdminUsersRequestObject) (GetAdminUsersResponseObject, error) {
	user := ctx.Value("user").(*auth.JWTClaims)
	if !user.IsAdmin {
		return GetAdminUsers403JSONResponse{
			Message: "Forbidden",
		}, nil
	}
	users, err := h.q.ListUsers(ctx)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	responseUsers := make([]User, len(users))
	for i, u := range users {
		responseUsers[i] = User{
			UserId:      int(u.UserID),
			Username:    u.Username,
			DisplayName: u.DisplayName,
			IconPath:    u.IconPath,
			IsAdmin:     u.IsAdmin,
		}
	}
	return GetAdminUsers200JSONResponse{
		Users: responseUsers,
	}, nil
}

func (h *ApiHandler) PostLogin(ctx context.Context, request PostLoginRequestObject) (PostLoginResponseObject, error) {
	username := request.Body.Username
	password := request.Body.Password
	userId, err := auth.Login(ctx, h.q, username, password)
	if err != nil {
		return PostLogin401JSONResponse{
			Message: "Invalid username or password",
		}, nil
	}

	user, err := h.q.GetUserById(ctx, int32(userId))
	if err != nil {
		return PostLogin401JSONResponse{
			Message: "Invalid username or password",
		}, nil
	}

	jwt, err := auth.NewJWT(&user)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return PostLogin200JSONResponse{
		Token: jwt,
	}, nil
}

func (h *ApiHandler) GetToken(ctx context.Context, request GetTokenRequestObject) (GetTokenResponseObject, error) {
	user := ctx.Value("user").(*auth.JWTClaims)
	newToken, err := auth.NewShortLivedJWT(user)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return GetToken200JSONResponse{
		Token: newToken,
	}, nil
}

func (h *ApiHandler) GetGames(ctx context.Context, request GetGamesRequestObject) (GetGamesResponseObject, error) {
	user := ctx.Value("user").(*auth.JWTClaims)
	playerId := request.Params.PlayerId
	if !user.IsAdmin {
		if playerId == nil || *playerId != user.UserID {
			return GetGames403JSONResponse{
				Message: "Forbidden",
			}, nil
		}
	}
	if playerId == nil {
		gameRows, err := h.q.ListGames(ctx)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		games := make([]Game, len(gameRows))
		for i, row := range gameRows {
			var startedAt *int
			if row.StartedAt.Valid {
				startedAtTimestamp := int(row.StartedAt.Time.Unix())
				startedAt = &startedAtTimestamp
			}
			var problem *Problem
			if row.ProblemID != nil {
				if row.Title == nil || row.Description == nil {
					panic("inconsistent data")
				}
				problem = &Problem{
					ProblemId:   int(*row.ProblemID),
					Title:       *row.Title,
					Description: *row.Description,
				}
			}
			games[i] = Game{
				GameId:          int(row.GameID),
				State:           GameState(row.State),
				DisplayName:     row.DisplayName,
				DurationSeconds: int(row.DurationSeconds),
				StartedAt:       startedAt,
				Problem:         problem,
			}
		}
		return GetGames200JSONResponse{
			Games: games,
		}, nil
	} else {
		gameRows, err := h.q.ListGamesForPlayer(ctx, int32(*playerId))
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		games := make([]Game, len(gameRows))
		for i, row := range gameRows {
			var startedAt *int
			if row.StartedAt.Valid {
				startedAtTimestamp := int(row.StartedAt.Time.Unix())
				startedAt = &startedAtTimestamp
			}
			var problem *Problem
			if row.ProblemID != nil {
				if row.Title == nil || row.Description == nil {
					panic("inconsistent data")
				}
				problem = &Problem{
					ProblemId:   int(*row.ProblemID),
					Title:       *row.Title,
					Description: *row.Description,
				}
			}
			games[i] = Game{
				GameId:          int(row.GameID),
				State:           GameState(row.State),
				DisplayName:     row.DisplayName,
				DurationSeconds: int(row.DurationSeconds),
				StartedAt:       startedAt,
				Problem:         problem,
			}
		}
		return GetGames200JSONResponse{
			Games: games,
		}, nil
	}
}

func (h *ApiHandler) GetGamesGameId(ctx context.Context, request GetGamesGameIdRequestObject) (GetGamesGameIdResponseObject, error) {
	user := ctx.Value("user").(*auth.JWTClaims)
	// TODO: check user permission
	gameId := request.GameId
	row, err := h.q.GetGameById(ctx, int32(gameId))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return GetGamesGameId404JSONResponse{
				Message: "Game not found",
			}, nil
		} else {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	var startedAt *int
	if row.StartedAt.Valid {
		startedAtTimestamp := int(row.StartedAt.Time.Unix())
		startedAt = &startedAtTimestamp
	}
	var problem *Problem
	if row.ProblemID != nil {
		if row.Title == nil || row.Description == nil {
			panic("inconsistent data")
		}
		if user.IsAdmin || (GameState(row.State) != GameStateClosed && GameState(row.State) != GameStateWaitingEntries) {
			problem = &Problem{
				ProblemId:   int(*row.ProblemID),
				Title:       *row.Title,
				Description: *row.Description,
			}
		}
	}
	game := Game{
		GameId:          int(row.GameID),
		State:           GameState(row.State),
		DisplayName:     row.DisplayName,
		DurationSeconds: int(row.DurationSeconds),
		StartedAt:       startedAt,
		Problem:         problem,
	}
	return GetGamesGameId200JSONResponse{
		Game: game,
	}, nil
}

func _assertUserResponseIsCompatibleWithJWTClaims() {
	var c auth.JWTClaims
	var u User
	u.UserId = c.UserID
	u.Username = c.Username
	u.DisplayName = c.DisplayName
	u.IconPath = c.IconPath
	u.IsAdmin = c.IsAdmin
	_ = u
}

func setupJWTFromAuthorizationHeader(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	const prefix = "Bearer "
	if !strings.HasPrefix(authorization, prefix) {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	token := authorization[len(prefix):]
	claims, err := auth.ParseJWT(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), "user", claims)))
	return nil
}

func NewJWTMiddleware() StrictMiddlewareFunc {
	return func(handler StrictHandlerFunc, operationID string) StrictHandlerFunc {
		if operationID == "PostLogin" {
			return handler
		}

		return func(c echo.Context, request interface{}) (interface{}, error) {
			err := setupJWTFromAuthorizationHeader(c)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}
			return handler(c, request)
		}
	}
}
