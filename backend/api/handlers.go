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
			if row.ProblemID.Valid {
				if !row.Title.Valid || !row.Description.Valid {
					panic("inconsistent data")
				}
				problem = &Problem{
					ProblemId:   int(row.ProblemID.Int32),
					Title:       row.Title.String,
					Description: row.Description.String,
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
			if row.ProblemID.Valid {
				if !row.Title.Valid || !row.Description.Valid {
					panic("inconsistent data")
				}
				problem = &Problem{
					ProblemId:   int(row.ProblemID.Int32),
					Title:       row.Title.String,
					Description: row.Description.String,
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
	// TODO: user permission
	// user := ctx.Value("user").(*auth.JWTClaims)
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
	if row.ProblemID.Valid && GameState(row.State) != Closed && GameState(row.State) != WaitingEntries {
		if !row.Title.Valid || !row.Description.Valid {
			panic("inconsistent data")
		}
		problem = &Problem{
			ProblemId:   int(row.ProblemID.Int32),
			Title:       row.Title.String,
			Description: row.Description.String,
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

func _assertJwtPayloadIsCompatibleWithJWTClaims() {
	var c auth.JWTClaims
	var p JwtPayload
	p.UserId = c.UserID
	p.Username = c.Username
	p.DisplayName = c.DisplayName
	p.IconPath = c.IconPath
	p.IsAdmin = c.IsAdmin
	_ = p
}

func NewJWTMiddleware() StrictMiddlewareFunc {
	return func(handler StrictHandlerFunc, operationID string) StrictHandlerFunc {
		if operationID == "PostLogin" {
			return handler
		} else {
			return func(c echo.Context, request interface{}) (response interface{}, err error) {
				authorization := c.Request().Header.Get("Authorization")
				const prefix = "Bearer "
				if !strings.HasPrefix(authorization, prefix) {
					return nil, echo.NewHTTPError(http.StatusUnauthorized)
				}
				token := authorization[len(prefix):]

				claims, err := auth.ParseJWT(token)
				if err != nil {
					return nil, echo.NewHTTPError(http.StatusUnauthorized)
				}
				c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), "user", claims)))
				return handler(c, request)
			}
		}
	}
}
