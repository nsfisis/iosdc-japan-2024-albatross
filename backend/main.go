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

	openAPISpec, err := api.GetSwaggerWithPrefix("/iosdc-japan/2024/code-battle/api")
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

	gameHubs := game.NewGameHubs(queries, taskQueue, workerServer.Results())
	err = gameHubs.RestoreFromDB(ctx)
	if err != nil {
		log.Fatalf("Error restoring game hubs from db %v", err)
	}
	defer gameHubs.Close()
	sockGroup := e.Group("/iosdc-japan/2024/code-battle/sock")
	sockHandler := gameHubs.SockHandler()
	sockGroup.GET("/golf/:gameID/play", func(c echo.Context) error {
		return sockHandler.HandleSockGolfPlay(c)
	})
	sockGroup.GET("/golf/:gameID/watch", func(c echo.Context) error {
		return sockHandler.HandleSockGolfWatch(c)
	})

	apiGroup := e.Group("/iosdc-japan/2024/code-battle/api")
	apiGroup.Use(oapimiddleware.OapiRequestValidator(openAPISpec))
	apiHandler := api.NewHandler(queries, gameHubs)
	api.RegisterHandlers(apiGroup, api.NewStrictHandler(apiHandler, nil))

	adminHandler := admin.NewHandler(queries, gameHubs)
	adminGroup := e.Group("/iosdc-japan/2024/code-battle/admin")
	adminHandler.RegisterHandlers(adminGroup)

	if config.isLocal {
		// For local dev: This is never used in production because the reverse
		// proxy directly handles /files.
		filesGroup := e.Group("/iosdc-japan/2024/code-battle/files")
		filesGroup.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Root:       "/",
			Filesystem: http.Dir("/data/files"),
			IgnoreBase: true,
		}))

		// For local dev: This is never used in production because the reverse
		// proxy sends these paths to the app server.
		e.GET("/iosdc-japan/2024/code-battle/*", func(c echo.Context) error {
			return c.Redirect(http.StatusPermanentRedirect, "http://localhost:5173"+c.Request().URL.Path)
		})
		e.POST("/iosdc-japan/2024/code-battle/*", func(c echo.Context) error {
			return c.Redirect(http.StatusPermanentRedirect, "http://localhost:5173"+c.Request().URL.Path)
		})
	}

	go gameHubs.Run()

	go func() {
		if err := workerServer.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := e.Start(":80"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
