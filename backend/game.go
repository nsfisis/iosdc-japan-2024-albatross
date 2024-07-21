package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type GameHub struct {
	game              *Game
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

func NewGameHub(game *Game) *GameHub {
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

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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

func serveWs(hub *GameHub, w http.ResponseWriter, r *http.Request, team string) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	client := &GameClient{hub: hub, conn: conn, send: make(chan *Message), team: team}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
	return nil
}

func serveWsWatcher(hub *GameHub, w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	watcher := &GameWatcher{hub: hub, conn: conn, send: make(chan *Message)}
	watcher.hub.registerWatcher <- watcher

	go watcher.writePump()
	go watcher.readPump()
	return nil
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
