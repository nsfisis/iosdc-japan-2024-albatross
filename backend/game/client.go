package game

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type playerClient struct {
	hub         *gameHub
	conn        *websocket.Conn
	s2cMessages chan playerMessageS2C
	playerID    int
}

// Receives messages from the client and sends them to the hub.
func (c *playerClient) readPump() {
	defer func() {
		log.Printf("closing player client")
		c.hub.unregisterPlayer <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var rawMessage map[string]json.RawMessage
		if err := c.conn.ReadJSON(&rawMessage); err != nil {
			log.Printf("error: %v", err)
			return
		}
		message, err := asPlayerMessageC2S(rawMessage)
		if err != nil {
			log.Printf("error: %v", err)
			return
		}
		c.hub.playerC2SMessages <- &playerMessageC2SWithClient{c, message}
	}
}

// Receives messages from the hub and sends them to the client.
func (c *playerClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.s2cMessages:
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

type watcherClient struct {
	hub         *gameHub
	conn        *websocket.Conn
	s2cMessages chan watcherMessageS2C
}

// Receives messages from the client and sends them to the hub.
func (c *watcherClient) readPump() {
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
}

// Receives messages from the hub and sends them to the client.
func (c *watcherClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
		log.Printf("closing watcher client")
		c.hub.unregisterWatcher <- c
	}()
	for {
		select {
		case message, ok := <-c.s2cMessages:
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
