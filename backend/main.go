package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	oapimiddleware "github.com/oapi-codegen/echo-middleware"

	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/api"
	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/db"
	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/game"
)

func connectDB(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}

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

	dbDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.dbHost, config.dbPort, config.dbUser, config.dbPassword, config.dbName)
	connPool, err := connectDB(ctx, dbDSN)
	if err != nil {
		log.Fatalf("Error connecting to db %v", err)
	}
	defer connPool.Close()

	queries := db.New(connPool)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	gameHubs := game.NewGameHubs(queries)
	err = gameHubs.RestoreFromDB(ctx)
	if err != nil {
		log.Fatalf("Error restoring game hubs from db %v", err)
	}
	defer gameHubs.Close()
	sockGroup := e.Group("/sock")
	sockHandler := gameHubs.SockHandler()
	sockGroup.GET("/golf/:gameId/play", func(c echo.Context) error {
		return sockHandler.HandleSockGolfPlay(c)
	})
	sockGroup.GET("/golf/:gameId/watch", func(c echo.Context) error {
		return sockHandler.HandleSockGolfWatch(c)
	})

	apiGroup := e.Group("/api")
	apiGroup.Use(oapimiddleware.OapiRequestValidator(openApiSpec))
	apiHandler := api.NewHandler(queries, gameHubs)
	api.RegisterHandlers(apiGroup, api.NewStrictHandler(apiHandler, []api.StrictMiddlewareFunc{
		api.NewJWTMiddleware(),
	}))

	gameHubs.Run()

	if err := e.Start(":80"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
