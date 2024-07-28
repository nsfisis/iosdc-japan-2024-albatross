package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	oapimiddleware "github.com/oapi-codegen/echo-middleware"

	"github.com/nsfisis/iosdc-2024-albatross/backend/api"
	"github.com/nsfisis/iosdc-2024-albatross/backend/db"
)

const (
	gameTypeGolf = "golf"
)

const (
	gameStateWaiting  = "waiting"
	gameStateReady    = "ready"
	gameStatePlaying  = "playing"
	gameStateFinished = "finished"
)

type Game struct {
	GameID int `db:"game_id"`
	// "golf"
	Type      string `db:"type"`
	CreatedAt string `db:"created_at"`
	State     string `db:"state"`
}

var gameHubs = map[int]*GameHub{}

func startGame(game *Game) {
	if gameHubs[game.GameID] != nil {
		return
	}
	gameHubs[game.GameID] = NewGameHub(game)
	go gameHubs[game.GameID].Run()
}

/*
func handleGolfPost(w http.ResponseWriter, r *http.Request) {
	var yourTeam string
	waitingGolfGames := []Game{}
	err := db.Select(&waitingGolfGames, "SELECT * FROM games WHERE type = $1 AND state = $2 ORDER BY created_at", gameTypeGolf, gameStateWaiting)
	if err != nil {
		http.Error(w, "Error getting games", http.StatusInternalServerError)
		return
	}
	if len(waitingGolfGames) == 0 {
		_, err = db.Exec("INSERT INTO games (type, state) VALUES ($1, $2)", gameTypeGolf, gameStateWaiting)
		if err != nil {
			http.Error(w, "Error creating game", http.StatusInternalServerError)
			return
		}
		waitingGolfGames = []Game{}
		err = db.Select(&waitingGolfGames, "SELECT * FROM games WHERE type = $1 AND state = $2 ORDER BY created_at", gameTypeGolf, gameStateWaiting)
		if err != nil {
			http.Error(w, "Error getting games", http.StatusInternalServerError)
			return
		}
		yourTeam = "a"
		startGame(&waitingGolfGames[0])
	} else {
		yourTeam = "b"
		db.Exec("UPDATE games SET state = $1 WHERE game_id = $2", gameStateReady, waitingGolfGames[0].GameID)
	}
	waitingGame := waitingGolfGames[0]

	http.Redirect(w, r, fmt.Sprintf("/golf/%d/%s/", waitingGame.GameID, yourTeam), http.StatusSeeOther)
}
*/

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

	{
		apiGroup := e.Group("/api")
		apiGroup.Use(oapimiddleware.OapiRequestValidator(openApiSpec))
		apiHandler := api.NewHandler(queries)
		api.RegisterHandlers(apiGroup, api.NewStrictHandler(apiHandler, []api.StrictMiddlewareFunc{
			api.NewJWTMiddleware(),
		}))
	}

	e.GET("/sock/golf/:gameId/watch", func(c echo.Context) error {
		gameId := c.Param("gameId")
		gameIdInt, err := strconv.Atoi(gameId)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid game id")
		}
		var hub *GameHub
		for _, h := range gameHubs {
			if h.game.GameID == gameIdInt {
				hub = h
				break
			}
		}
		if hub == nil {
			return echo.NewHTTPError(http.StatusNotFound, "Game not found")
		}
		return serveWsWatcher(hub, c.Response(), c.Request())
	})

	e.GET("/sock/golf/:gameId/play", func(c echo.Context) error {
		gameId := c.Param("gameId")
		gameIdInt, err := strconv.Atoi(gameId)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid game id")
		}
		var hub *GameHub
		for _, h := range gameHubs {
			if h.game.GameID == gameIdInt {
				hub = h
				break
			}
		}
		if hub == nil {
			return echo.NewHTTPError(http.StatusNotFound, "Game not found")
		}
		return serveWs(hub, c.Response(), c.Request(), "a")
	})

	defer func() {
		for _, hub := range gameHubs {
			hub.Close()
		}
	}()

	if err := e.Start(":80"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
