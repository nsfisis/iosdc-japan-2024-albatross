package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/nsfisis/iosdc-2024-albatross/backend/auth"
	"github.com/nsfisis/iosdc-2024-albatross/backend/db"
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
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
	return GetGamesGameId200JSONResponse(game), nil
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
