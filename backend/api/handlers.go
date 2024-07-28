package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/nsfisis/iosdc-2024-albatross-backend/auth"
	"github.com/nsfisis/iosdc-2024-albatross-backend/db"
)

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
		}, echo.NewHTTPError(http.StatusUnauthorized, "Invalid username or password")
	}

	user, err := h.q.GetUserById(ctx, int32(userId))
	if err != nil {
		return PostLogin401JSONResponse{
			Message: "Invalid username or password",
		}, echo.NewHTTPError(http.StatusUnauthorized, "Invalid username or password")
	}

	jwt, err := auth.NewJWT(&user)
	if err != nil {
		// TODO
		return PostLogin401JSONResponse{
			Message: "Internal Server Error",
		}, echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	return PostLogin200JSONResponse{
		Token: jwt,
	}, nil
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
