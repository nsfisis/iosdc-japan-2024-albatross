package game

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: insecure!
		return true
	},
}

func servePlayerWs(hub *gameHub, w http.ResponseWriter, r *http.Request, playerID int) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	player := &playerClient{
		hub:         hub,
		conn:        conn,
		s2cMessages: make(chan playerMessageS2C),
		playerID:    playerID,
	}
	hub.registerPlayer <- player

	go player.writePump()
	go player.readPump()
	return nil
}

func serveWatcherWs(hub *gameHub, w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	watcher := &watcherClient{
		hub:         hub,
		conn:        conn,
		s2cMessages: make(chan watcherMessageS2C),
	}
	hub.registerWatcher <- watcher

	go watcher.writePump()
	go watcher.readPump()
	return nil
}
