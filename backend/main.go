package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"

	"github.com/nsfisis/iosdc-2024-albatross-backend/auth"
	"github.com/nsfisis/iosdc-2024-albatross-backend/db"
)

type Config struct {
	dbHost     string
	dbPort     string
	dbUser     string
	dbPassword string
	dbName     string
}

var config *Config

func loadEnv() (*Config, error) {
	dbHost, exists := os.LookupEnv("ALBATROSS_DB_HOST")
	if !exists {
		return nil, fmt.Errorf("ALBATROSS_DB_HOST not set")
	}
	dbPort, exists := os.LookupEnv("ALBATROSS_DB_PORT")
	if !exists {
		return nil, fmt.Errorf("ALBATROSS_DB_PORT not set")
	}
	dbUser, exists := os.LookupEnv("ALBATROSS_DB_USER")
	if !exists {
		return nil, fmt.Errorf("ALBATROSS_DB_USER not set")
	}
	dbPassword, exists := os.LookupEnv("ALBATROSS_DB_PASSWORD")
	if !exists {
		return nil, fmt.Errorf("ALBATROSS_DB_PASSWORD not set")
	}
	dbName, exists := os.LookupEnv("ALBATROSS_DB_NAME")
	if !exists {
		return nil, fmt.Errorf("ALBATROSS_DB_NAME not set")
	}
	return &Config{
		dbHost:     dbHost,
		dbPort:     dbPort,
		dbUser:     dbUser,
		dbPassword: dbPassword,
		dbName:     dbName,
	}, nil
}

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

func handleApiLogin(c echo.Context, queries *db.Queries) error {
	type LoginRequestData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type LoginResponseData struct {
		Token string `json:"token"`
	}

	ctx := c.Request().Context()

	requestData := new(LoginRequestData)
	if err := c.Bind(requestData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userId, err := auth.Login(ctx, queries, requestData.Username, requestData.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	user, err := queries.GetUserById(ctx, int32(userId))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	jwt, err := auth.NewJWT(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	responseData := LoginResponseData{
		Token: jwt,
	}

	return c.JSON(http.StatusOK, responseData)
}

func main() {
	var err error
	config, err = loadEnv()
	if err != nil {
		fmt.Printf("Error loading env %v", err)
		return
	}

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.dbHost, config.dbPort, config.dbUser, config.dbPassword, config.dbName))
	if err != nil {
		log.Fatalf("Error connecting to db %v", err)
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	e := echo.New()

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

	e.POST("/api/login", func(c echo.Context) error {
		return handleApiLogin(c, queries)
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
