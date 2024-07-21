package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"

	"iosdc-code-battle-poc/db"
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

func handleGolfPost(w http.ResponseWriter, r *http.Request) {
	/*
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
	*/
}

func handleApiLogin(w http.ResponseWriter, r *http.Request, queries *db.Queries) {
	type LoginRequestData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type LoginResponseData struct {
		UserId int `json:"userId"`
	}

	ctx := r.Context()

	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		http.Error(w, "Content-Type is not application/json", http.StatusBadRequest)
		return
	}

	var requestData LoginRequestData
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestData)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	userId, err := authLogin(ctx, queries, requestData.Username, requestData.Password)
	if err != nil {
		http.Error(w, "Login failed", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	responseData := LoginResponseData{
		UserId: userId,
	}
	encoder := json.NewEncoder(w)
	err = encoder.Encode(responseData)
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
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

	server := http.NewServeMux()

	server.HandleFunc("GET /assets/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public"+r.URL.Path)
	})

	server.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/index.html")
	})

	server.HandleFunc("POST /golf/{$}", func(w http.ResponseWriter, r *http.Request) {
		handleGolfPost(w, r)
	})

	server.HandleFunc("GET /sock/golf/{gameId}/watch/{$}", func(w http.ResponseWriter, r *http.Request) {
		gameId := r.PathValue("gameId")
		gameIdInt, err := strconv.Atoi(gameId)
		if err != nil {
			http.Error(w, "Invalid game id", http.StatusBadRequest)
			return
		}
		var hub *GameHub
		for _, h := range gameHubs {
			if h.game.GameID == gameIdInt {
				hub = h
				break
			}
		}
		if hub == nil {
			http.Error(w, "Game not found", http.StatusNotFound)
			return
		}
		serveWsWatcher(hub, w, r)
	})

	server.HandleFunc("GET /sock/golf/{gameId}/{team}/{$}", func(w http.ResponseWriter, r *http.Request) {
		gameId := r.PathValue("gameId")
		gameIdInt, err := strconv.Atoi(gameId)
		if err != nil {
			http.Error(w, "Invalid game id", http.StatusBadRequest)
			return
		}
		var hub *GameHub
		for _, h := range gameHubs {
			if h.game.GameID == gameIdInt {
				hub = h
				break
			}
		}
		if hub == nil {
			http.Error(w, "Game not found", http.StatusNotFound)
			return
		}
		team := r.PathValue("team")
		serveWs(hub, w, r, team)
	})

	server.HandleFunc("POST /api/login", func(w http.ResponseWriter, r *http.Request) {
		handleApiLogin(w, r, queries)
	})

	defer func() {
		for _, hub := range gameHubs {
			hub.Close()
		}
	}()

	http.ListenAndServe(":80", server)
}
