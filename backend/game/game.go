package game

import (
	"context"
	"log"
	"time"

	"github.com/gorilla/websocket"

	"github.com/nsfisis/iosdc-2024-albatross/backend/api"
	"github.com/nsfisis/iosdc-2024-albatross/backend/db"
)

type gameState = api.GameState

const (
	gameStateClosed         gameState = api.Closed
	gameStateWaitingEntries gameState = api.WaitingEntries
	gameStateWaitingStart   gameState = api.WaitingStart
	gameStatePrepare        gameState = api.Prepare
	gameStateStarting       gameState = api.Starting
	gameStateGaming         gameState = api.Gaming
	gameStateFinished       gameState = api.Finished
)

type game struct {
	gameID          int
	state           string
	displayName     string
	durationSeconds int
	startedAt       *time.Time
	problem         *problem
}

type problem struct {
	problemID   int
	title       string
	description string
}

// func startGame(game *Game) {
// 	if gameHubs[game.GameID] != nil {
// 		return
// 	}
// 	gameHubs[game.GameID] = NewGameHub(game)
// 	go gameHubs[game.GameID].Run()
// }

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

type GameHub struct {
	game              *game
	clients           map[*GameClient]bool
	receive           chan *MessageWithClient
	register          chan *GameClient
	unregister        chan *GameClient
	watchers          map[*GameWatcher]bool
	registerWatcher   chan *GameWatcher
	unregisterWatcher chan *GameWatcher
	state             int
	finishTime        time.Time
}

func NewGameHub(game *game) *GameHub {
	return &GameHub{
		game:              game,
		clients:           make(map[*GameClient]bool),
		receive:           make(chan *MessageWithClient),
		register:          make(chan *GameClient),
		unregister:        make(chan *GameClient),
		watchers:          make(map[*GameWatcher]bool),
		registerWatcher:   make(chan *GameWatcher),
		unregisterWatcher: make(chan *GameWatcher),
		state:             0,
	}
}

func (h *GameHub) Run() {
	ticker := time.NewTicker(10 * time.Second)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("client registered: %d", len(h.clients))
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				h.closeClient(client)
			}
			log.Printf("client unregistered: %d", len(h.clients))
			if len(h.clients) == 0 {
				h.Close()
				return
			}
		case watcher := <-h.registerWatcher:
			h.watchers[watcher] = true
			log.Printf("watcher registered: %d", len(h.watchers))
		case watcher := <-h.unregisterWatcher:
			if _, ok := h.watchers[watcher]; ok {
				h.closeWatcher(watcher)
			}
			log.Printf("watcher unregistered: %d", len(h.watchers))
		case message := <-h.receive:
			log.Printf("received message: %s", message.Message.Type)
			switch message.Message.Type {
			case "connect":
				if h.state == 0 {
					h.state = 1
				} else if h.state == 1 {
					h.state = 2
					for client := range h.clients {
						client.send <- &Message{Type: "prepare", Data: MessageDataPrepare{Problem: "1 から 100 までの FizzBuzz を実装せよ (終端を含む)。"}}
					}
				} else {
					log.Printf("invalid state: %d", h.state)
					h.closeClient(message.Client)
				}
			case "ready":
				if h.state == 2 {
					h.state = 3
				} else if h.state == 3 {
					h.state = 4
					for client := range h.clients {
						client.send <- &Message{Type: "start", Data: MessageDataStart{StartTime: time.Now().Add(10 * time.Second).UTC().Format(time.RFC3339)}}
					}
					h.finishTime = time.Now().Add(3 * time.Minute)
				} else {
					log.Printf("invalid state: %d", h.state)
					h.closeClient(message.Client)
				}
			case "code":
				if h.state == 4 {
					code := message.Message.Data.(MessageDataCode).Code
					message.Client.code = code
					message.Client.send <- &Message{Type: "score", Data: MessageDataScore{Score: 100}}
					if message.Client.score == nil {
						message.Client.score = new(int)
					}
					*message.Client.score = 100

					var scoreA, scoreB *int
					var codeA, codeB string
					for client := range h.clients {
						if client.team == "a" {
							scoreA = client.score
							codeA = client.code
						} else {
							scoreB = client.score
							codeB = client.code
						}
					}
					for watcher := range h.watchers {
						watcher.send <- &Message{
							Type: "watch",
							Data: MessageDataWatch{
								Problem: "1 から 100 までの FizzBuzz を実装せよ (終端を含む)。",
								ScoreA:  scoreA,
								CodeA:   codeA,
								ScoreB:  scoreB,
								CodeB:   codeB,
							},
						}
					}
				} else {
					log.Printf("invalid state: %d", h.state)
					h.closeClient(message.Client)
				}
			default:
				log.Printf("unknown message type: %s", message.Message.Type)
				h.closeClient(message.Client)
			}
		case <-ticker.C:
			log.Printf("state: %d", h.state)
			if h.state == 4 {
				if time.Now().After(h.finishTime) {
					h.state = 5
					clientAndScores := make(map[*GameClient]*int)
					for client := range h.clients {
						clientAndScores[client] = client.score
					}
					for client, score := range clientAndScores {
						var opponentScore *int
						for c2, s2 := range clientAndScores {
							if c2 != client {
								opponentScore = s2
								break
							}
						}
						client.send <- &Message{Type: "finish", Data: MessageDataFinish{YourScore: score, OpponentScore: opponentScore}}
					}
				}
			}
		}
	}
}

func (h *GameHub) Close() {
	for client := range h.clients {
		h.closeClient(client)
	}
	close(h.receive)
	close(h.register)
	close(h.unregister)
	for watcher := range h.watchers {
		h.closeWatcher(watcher)
	}
	close(h.registerWatcher)
	close(h.unregisterWatcher)
}

func (h *GameHub) closeClient(client *GameClient) {
	delete(h.clients, client)
	close(client.send)
}

func (h *GameHub) closeWatcher(watcher *GameWatcher) {
	delete(h.watchers, watcher)
	close(watcher.send)
}

type GameClient struct {
	hub   *GameHub
	conn  *websocket.Conn
	send  chan *Message
	score *int
	code  string
	team  string
}

type GameWatcher struct {
	hub  *GameHub
	conn *websocket.Conn
	send chan *Message
}

// Receives messages from the client and sends them to the hub.
func (c *GameClient) readPump() {
	defer func() {
		log.Printf("closing client")
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var message Message
		err := c.conn.ReadJSON(&message)
		if err != nil {
			log.Printf("error: %v", err)
			return
		}
		c.hub.receive <- &MessageWithClient{c, &message}
	}
}

// Receives messages from the hub and sends them to the client.
func (c *GameClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteJSON(message)
			if err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Receives messages from the client and sends them to the hub.
func (c *GameWatcher) readPump() {
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
}

// Receives messages from the hub and sends them to the client.
func (c *GameWatcher) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
		log.Printf("closing watcher")
		c.hub.unregisterWatcher <- c
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteJSON(message)
			if err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

type GameHubs struct {
	hubs map[int]*GameHub
}

func NewGameHubs() *GameHubs {
	return &GameHubs{
		hubs: make(map[int]*GameHub),
	}
}

func (hubs *GameHubs) Close() {
	for _, hub := range hubs.hubs {
		hub.Close()
	}
}

func (hubs *GameHubs) RestoreFromDB(ctx context.Context, q *db.Queries) error {
	games, err := q.ListGames(ctx)
	if err != nil {
		return err
	}
	_ = games
	return nil
}

func (hubs *GameHubs) SockHandler() *sockHandler {
	return newSockHandler(hubs)
}
