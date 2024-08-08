package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"

	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/auth"
	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/db"
)

type APIHandler struct {
	q    *db.Queries
	hubs GameHubsInterface
}

type GameHubsInterface interface {
	StartGame(gameID int) error
}

func (h *APIHandler) PostLogin(ctx context.Context, request PostLoginRequestObject) (PostLoginResponseObject, error) {
	username := request.Body.Username
	password := request.Body.Password
	userID, err := auth.Login(ctx, h.q, username, password)
	if err != nil {
		return PostLogin401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{
				Message: "Invalid username or password",
			},
		}, nil
	}

	user, err := h.q.GetUserByID(ctx, int32(userID))
	if err != nil {
		return PostLogin401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{
				Message: "Invalid username or password",
			},
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

func (h *APIHandler) GetToken(ctx context.Context, request GetTokenRequestObject, user *auth.JWTClaims) (GetTokenResponseObject, error) {
	newToken, err := auth.NewShortLivedJWT(user)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return GetToken200JSONResponse{
		Token: newToken,
	}, nil
}

func (h *APIHandler) GetGames(ctx context.Context, request GetGamesRequestObject, user *auth.JWTClaims) (GetGamesResponseObject, error) {
	gameRows, err := h.q.ListGamesForPlayer(ctx, int32(user.UserID))
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
				ProblemID:   int(*row.ProblemID),
				Title:       *row.Title,
				Description: *row.Description,
			}
		}
		games[i] = Game{
			GameID:          int(row.GameID),
			GameType:        GameGameType(row.GameType),
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

func (h *APIHandler) GetGame(ctx context.Context, request GetGameRequestObject, user *auth.JWTClaims) (GetGameResponseObject, error) {
	// TODO: check user permission
	gameID := request.GameID
	row, err := h.q.GetGameByID(ctx, int32(gameID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return GetGame404JSONResponse{
				NotFoundJSONResponse: NotFoundJSONResponse{
					Message: "Game not found",
				},
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
		if user.IsAdmin || (GameState(row.State) != Closed && GameState(row.State) != WaitingEntries) {
			problem = &Problem{
				ProblemID:   int(*row.ProblemID),
				Title:       *row.Title,
				Description: *row.Description,
			}
		}
	}
	game := Game{
		GameID:          int(row.GameID),
		GameType:        GameGameType(row.GameType),
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
