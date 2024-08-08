package game

import (
	"encoding/json"
	"fmt"
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
func (c *playerClient) readPump() error {
	defer func() {
		log.Printf("closing player client")
		c.hub.unregisterPlayer <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		return err
	}
	c.conn.SetPongHandler(func(string) error { return c.conn.SetReadDeadline(time.Now().Add(pongWait)) })
	for {
		var rawMessage map[string]json.RawMessage
		if err := c.conn.ReadJSON(&rawMessage); err != nil {
			return err
		}
		message, err := asPlayerMessageC2S(rawMessage)
		if err != nil {
			return err
		}
		c.hub.playerC2SMessages <- &playerMessageC2SWithClient{c, message}
	}
}

// Receives messages from the hub and sends them to the client.
func (c *playerClient) writePump() error {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.s2cMessages:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return err
			}
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					return err
				}
				return fmt.Errorf("closing player client")
			}

			err := c.conn.WriteJSON(message)
			if err != nil {
				return err
			}
		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return err
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return err
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
func (c *watcherClient) readPump() error {
	c.conn.SetReadLimit(maxMessageSize)
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		return err
	}
	c.conn.SetPongHandler(func(string) error { return c.conn.SetReadDeadline(time.Now().Add(pongWait)) })
	return nil
}

// Receives messages from the hub and sends them to the client.
func (c *watcherClient) writePump() error {
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
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return err
			}
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					return err
				}
				return fmt.Errorf("closing watcher client")
			}

			err := c.conn.WriteJSON(message)
			if err != nil {
				return err
			}
		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return err
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return err
			}
		}
	}
}
