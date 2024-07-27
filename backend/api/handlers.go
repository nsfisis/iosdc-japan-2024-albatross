package api

import (
	"context"
	"net/http"

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

func (h *ApiHandler) PostApiLogin(ctx context.Context, request PostApiLoginRequestObject) (PostApiLoginResponseObject, error) {
	username := request.Body.Username
	password := request.Body.Password
	userId, err := auth.Login(ctx, h.q, username, password)
	if err != nil {
		return PostApiLogin401JSONResponse{
			Message: "Invalid username or password",
		}, echo.NewHTTPError(http.StatusUnauthorized, "Invalid username or password")
	}

	user, err := h.q.GetUserById(ctx, int32(userId))
	if err != nil {
		return PostApiLogin401JSONResponse{
			Message: "Invalid username or password",
		}, echo.NewHTTPError(http.StatusUnauthorized, "Invalid username or password")
	}

	jwt, err := auth.NewJWT(&user)
	if err != nil {
		// TODO
		return PostApiLogin401JSONResponse{
			Message: "Internal Server Error",
		}, echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	return PostApiLogin200JSONResponse{
		Token: jwt,
	}, nil
}

func _assertJwtPayloadIsCompatibleWithJWTClaims() {
	var c auth.JWTClaims
	var p JwtPayload
	p.UserId = float32(c.UserID)
	p.Username = c.Username
	p.DisplayUsername = c.DisplayUsername
	p.IconPath = c.IconPath
	p.IsAdmin = c.IsAdmin
	_ = p
}
