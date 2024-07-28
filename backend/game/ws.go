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
}

func servePlayerWs(hub *GameHub, w http.ResponseWriter, r *http.Request, team string) error {
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

func serveWatcherWs(hub *GameHub, w http.ResponseWriter, r *http.Request) error {
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
