package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"

	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/auth"
	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/db"
)

type ApiHandler struct {
	q    *db.Queries
	hubs GameHubsInterface
}

type GameHubsInterface interface {
	StartGame(gameID int) error
}

func (h *ApiHandler) AdminGetGames(ctx context.Context, request AdminGetGamesRequestObject, user *auth.JWTClaims) (AdminGetGamesResponseObject, error) {
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
	return AdminGetGames200JSONResponse{
		Games: games,
	}, nil
}

func (h *ApiHandler) AdminGetGame(ctx context.Context, request AdminGetGameRequestObject, user *auth.JWTClaims) (AdminGetGameResponseObject, error) {
	gameId := request.GameId
	row, err := h.q.GetGameById(ctx, int32(gameId))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return AdminGetGame404JSONResponse{
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
	return AdminGetGame200JSONResponse{
		Game: game,
	}, nil
}

func (h *ApiHandler) AdminPutGame(ctx context.Context, request AdminPutGameRequestObject, user *auth.JWTClaims) (AdminPutGameResponseObject, error) {
	gameID := request.GameId
	displayName := request.Body.DisplayName
	durationSeconds := request.Body.DurationSeconds
	problemID := request.Body.ProblemId
	startedAt := request.Body.StartedAt
	state := request.Body.State

	game, err := h.q.GetGameById(ctx, int32(gameID))
	if err != nil {
		if err == pgx.ErrNoRows {
			return AdminPutGame404JSONResponse{
				Message: "Game not found",
			}, nil
		} else {
			return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	var changedState string
	if state != nil {
		changedState = string(*state)
		// TODO:
		if changedState != game.State && changedState == "prepare" {
			h.hubs.StartGame(int(gameID))
		}
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
		return AdminPutGame400JSONResponse{
			Message: err.Error(),
		}, nil
	}

	return AdminPutGame204Response{}, nil
}

func (h *ApiHandler) AdminGetUsers(ctx context.Context, request AdminGetUsersRequestObject, user *auth.JWTClaims) (AdminGetUsersResponseObject, error) {
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
	return AdminGetUsers200JSONResponse{
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

func (h *ApiHandler) GetToken(ctx context.Context, request GetTokenRequestObject, user *auth.JWTClaims) (GetTokenResponseObject, error) {
	newToken, err := auth.NewShortLivedJWT(user)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return GetToken200JSONResponse{
		Token: newToken,
	}, nil
}

func (h *ApiHandler) GetGames(ctx context.Context, request GetGamesRequestObject, user *auth.JWTClaims) (GetGamesResponseObject, error) {
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

func (h *ApiHandler) GetGame(ctx context.Context, request GetGameRequestObject, user *auth.JWTClaims) (GetGameResponseObject, error) {
	// TODO: check user permission
	gameId := request.GameId
	row, err := h.q.GetGameById(ctx, int32(gameId))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return GetGame404JSONResponse{
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
	return GetGame200JSONResponse{
		Game: game,
	}, nil
}
