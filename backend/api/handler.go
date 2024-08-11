package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/nullable"

	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/auth"
	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/db"
)

type Handler struct {
	q    *db.Queries
	hubs GameHubsInterface
}

type GameHubsInterface interface {
	StartGame(gameID int) error
}

func (h *Handler) PostLogin(ctx context.Context, request PostLoginRequestObject) (PostLoginResponseObject, error) {
	username := request.Body.Username
	password := request.Body.Password
	registrationToken := request.Body.RegistrationToken
	userID, err := auth.Login(ctx, h.q, username, password, registrationToken)
	if err != nil {
		log.Printf("login failed: %v", err)
		var msg string
		if errors.Is(err, auth.ErrInvalidRegistrationToken) {
			msg = "登録用 URL が無効です。イベントスタッフにお声がけください"
		} else if errors.Is(err, auth.ErrNoRegistrationToken) {
			msg = "登録用 URL からログインしてください。登録用 URL は Connpass のイベントページに記載しています"
		} else if errors.Is(err, auth.ErrForteeLoginTimeout) {
			msg = "ログインに失敗しました"
		} else {
			msg = "ユーザー名またはパスワードが誤っています"
		}
		return PostLogin401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{
				Message: msg,
			},
		}, nil
	}

	user, err := h.q.GetUserByID(ctx, int32(userID))
	if err != nil {
		return PostLogin401JSONResponse{
			UnauthorizedJSONResponse: UnauthorizedJSONResponse{
				Message: "ログインに失敗しました",
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

func (h *Handler) GetToken(_ context.Context, _ GetTokenRequestObject, user *auth.JWTClaims) (GetTokenResponseObject, error) {
	newToken, err := auth.NewShortLivedJWT(user)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return GetToken200JSONResponse{
		Token: newToken,
	}, nil
}

func (h *Handler) GetGames(ctx context.Context, _ GetGamesRequestObject, user *auth.JWTClaims) (GetGamesResponseObject, error) {
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
		games[i] = Game{
			GameID:          int(row.GameID),
			GameType:        GameGameType(row.GameType),
			State:           GameState(row.State),
			DisplayName:     row.DisplayName,
			DurationSeconds: int(row.DurationSeconds),
			StartedAt:       startedAt,
			Problem: Problem{
				ProblemID:   int(row.ProblemID),
				Title:       row.Title,
				Description: row.Description,
			},
		}
	}
	return GetGames200JSONResponse{
		Games: games,
	}, nil
}

func (h *Handler) GetGame(ctx context.Context, request GetGameRequestObject, user *auth.JWTClaims) (GetGameResponseObject, error) {
	// TODO: check user permission
	_ = user
	gameID := request.GameID
	row, err := h.q.GetGameByID(ctx, int32(gameID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return GetGame404JSONResponse{
				NotFoundJSONResponse: NotFoundJSONResponse{
					Message: "Game not found",
				},
			}, nil
		}
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	var startedAt *int
	if row.StartedAt.Valid {
		startedAtTimestamp := int(row.StartedAt.Time.Unix())
		startedAt = &startedAtTimestamp
	}
	playerRows, err := h.q.ListGamePlayers(ctx, int32(gameID))
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	players := make([]User, len(playerRows))
	for i, playerRow := range playerRows {
		players[i] = User{
			UserID:      int(playerRow.UserID),
			Username:    playerRow.Username,
			DisplayName: playerRow.DisplayName,
			IconPath:    playerRow.IconPath,
			IsAdmin:     playerRow.IsAdmin,
		}
	}
	testcaseIDs, err := h.q.ListTestcaseIDsByGameID(ctx, int32(gameID))
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	verificationSteps := make([]VerificationStep, len(testcaseIDs)+1)
	verificationSteps[0] = VerificationStep{
		Label: "Compile",
	}
	for i, testcaseID := range testcaseIDs {
		verificationSteps[i+1] = VerificationStep{
			TestcaseID: nullable.NewNullableWithValue(int(testcaseID)),
			Label:      fmt.Sprintf("Testcase %d", i+1),
		}
	}
	game := Game{
		GameID:          int(row.GameID),
		GameType:        GameGameType(row.GameType),
		State:           GameState(row.State),
		DisplayName:     row.DisplayName,
		DurationSeconds: int(row.DurationSeconds),
		StartedAt:       startedAt,
		Problem: Problem{
			ProblemID:   int(row.ProblemID),
			Title:       row.Title,
			Description: row.Description,
		},
		Players:           players,
		VerificationSteps: verificationSteps,
	}
	return GetGame200JSONResponse{
		Game: game,
	}, nil
}
