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

	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/admin"
	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/api"
	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/db"
	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/game"
	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/taskqueue"
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
	e.Renderer = admin.NewRenderer()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	taskQueue := taskqueue.NewQueue("task-db:6379")
	workerServer := taskqueue.NewWorkerServer("task-db:6379", queries)

	gameHubs := game.NewGameHubs(queries, taskQueue)
	err = gameHubs.RestoreFromDB(ctx)
	if err != nil {
		log.Fatalf("Error restoring game hubs from db %v", err)
	}
	defer gameHubs.Close()
	sockGroup := e.Group("/sock")
	sockHandler := gameHubs.SockHandler()
	sockGroup.GET("/golf/:gameID/play", func(c echo.Context) error {
		return sockHandler.HandleSockGolfPlay(c)
	})
	sockGroup.GET("/golf/:gameID/watch", func(c echo.Context) error {
		return sockHandler.HandleSockGolfWatch(c)
	})

	apiGroup := e.Group("/api")
	apiGroup.Use(oapimiddleware.OapiRequestValidator(openApiSpec))
	apiHandler := api.NewHandler(queries, gameHubs)
	api.RegisterHandlers(apiGroup, api.NewStrictHandler(apiHandler, nil))

	adminHandler := admin.NewAdminHandler(queries, gameHubs)
	adminGroup := e.Group("/admin")
	adminHandler.RegisterHandlers(adminGroup)

	// For local dev: This is never used in production because the reverse
	// proxy sends /login and /logout to the app server.
	e.GET("/login", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "http://localhost:5173/login")
	})
	e.POST("/logout", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "http://localhost:5173/logout")
	})

	gameHubs.Run()

	go func() {
		workerServer.Run()
	}()

	if err := e.Start(":80"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
