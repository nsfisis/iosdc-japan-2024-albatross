package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	dbHost     string
	dbPort     string
	dbUser     string
	dbPassword string
	dbName     string
}

var config *Config

var db *sqlx.DB

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
	gameTypeRace = "race"
)

const (
	gameStateWaiting  = "waiting"
	gameStateReady    = "ready"
	gameStatePlaying  = "playing"
	gameStateFinished = "finished"
)

type Game struct {
	GameID int `db:"game_id"`
	// "golf" or "race"
	Type      string `db:"type"`
	CreatedAt string `db:"created_at"`
	State     string `db:"state"`
}

func initDB() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS games (
			game_id SERIAL PRIMARY KEY,
			type VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			state VARCHAR(255) NOT NULL
		);
	`)
	return err
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

func handleRacePost(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/race/1/a/", http.StatusSeeOther)
}

func main() {
	var err error
	config, err = loadEnv()
	if err != nil {
		fmt.Printf("Error loading env %v", err)
		return
	}

	for i := 0; i < 5; i++ {
		db, err = sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.dbHost, config.dbPort, config.dbUser, config.dbPassword, config.dbName))
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatalf("Error connecting to db %v", err)
	}
	defer db.Close()

	err = initDB()
	if err != nil {
		log.Fatalf("Error initializing db %v", err)
	}

	server := http.NewServeMux()

	server.HandleFunc("GET /assets/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public"+r.URL.Path)
	})

	server.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/index.html")
	})

	server.HandleFunc("GET /golf/{$}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/index.html")
	})

	server.HandleFunc("POST /golf/{$}", func(w http.ResponseWriter, r *http.Request) {
		handleGolfPost(w, r)
	})

	server.HandleFunc("GET /golf/{gameId}/watch/{$}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/index.html")
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

	server.HandleFunc("GET /golf/{gameId}/{team}/{$}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/index.html")
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

	defer func() {
		for _, hub := range gameHubs {
			hub.Close()
		}
	}()

	http.ListenAndServe(":80", server)
}
