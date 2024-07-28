package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	oapimiddleware "github.com/oapi-codegen/echo-middleware"

	"github.com/nsfisis/iosdc-2024-albatross/backend/api"
	"github.com/nsfisis/iosdc-2024-albatross/backend/db"
	"github.com/nsfisis/iosdc-2024-albatross/backend/game"
)

func main() {
	var err error
	config, err := NewConfigFromEnv()
	if err != nil {
		log.Fatalf("Error loading env %v", err)
	}

	openApiSpec, err := api.GetSwaggerWithPrefix("/api")
	if err != nil {
		log.Fatalf("Error loading OpenAPI spec\n: %s", err)
	}

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.dbHost, config.dbPort, config.dbUser, config.dbPassword, config.dbName))
	if err != nil {
		log.Fatalf("Error connecting to db %v", err)
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	apiGroup := e.Group("/api")
	apiGroup.Use(oapimiddleware.OapiRequestValidator(openApiSpec))
	apiHandler := api.NewHandler(queries)
	api.RegisterHandlers(apiGroup, api.NewStrictHandler(apiHandler, []api.StrictMiddlewareFunc{
		api.NewJWTMiddleware(),
	}))

	gameHubs := game.NewGameHubs()
	err = gameHubs.RestoreFromDB(ctx, queries)
	if err != nil {
		log.Fatalf("Error restoring game hubs from db %v", err)
	}
	defer gameHubs.Close()
	sockGroup := e.Group("/sock")
	sockHandler := gameHubs.SockHandler()
	sockGroup.GET("/golf/:gameId/watch", func(c echo.Context) error {
		return sockHandler.HandleSockGolfWatch(c)
	})
	sockGroup.GET("/golf/:gameId/play", func(c echo.Context) error {
		return sockHandler.HandleSockGolfPlay(c)
	})

	if err := e.Start(":80"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
