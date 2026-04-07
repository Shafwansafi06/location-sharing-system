package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Maximum amount of time to push a frame out before failing
	writeWait = 10 * time.Second
	// Wait time before giving up on hearing back from client
	pongWait = 60 * time.Second
	// Period automated pings are dispatched (slightly under pongWait)
	pingPeriod = (pongWait * 9) / 10
	// Hard payload max byte limit defending against intentional memory flood attacks
	maxMsgSize = 1024
)

// Client acts as a middleman interfacing the actual WebSocket layer to the abstract Hub.
type Client struct {
	hub     *Hub            // Pointer to global Hub instance
	conn    *websocket.Conn // Upgraded websocket connection
	send    chan []byte     // Buffered queue of messages intended specifically for this client
	GroupID string          // Which room they pertain to
	UserID  string          // Unique identifier (UUID/username)
}

// ReadPump continuously pulls incoming messages from the frontend browser.
// It funnels these JSON texts towards the `hub.broadcast` safely.
func (c *Client) ReadPump() {
	// Guaranteed cleanup execution whenever this function exits.
	defer func() {
		c.hub.unregister <- c // Tells Hub to delete us from mapping
		c.conn.Close()        // Forces connection socket shutdown
	}()

	c.conn.SetReadLimit(maxMsgSize)
	// Initialize the deadline
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	// Setup Pong interceptor handler (refreshes deadline on every successful Pong)
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		// Wait infinitely for the next message from the internet...
		_, payload, err := c.conn.ReadMessage()
		if err != nil {
			// Ignore basic closing protocols, report actual network failures.
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break // Kill the ReadPump (triggers Defer Cleanup)
		}

		// Immediately inject new coordinates towards Hub controller.
		c.hub.broadcast <- HubMessage{Sender: c, Payload: payload}
	}
}

// WritePump continuously drains the `send` channel outputting bytes sequentially down to the browser.
func (c *Client) WritePump() {
	// Launch an automated trigger handling heartbeat
	ticker := time.NewTicker(pingPeriod)

	// Ensure guaranteed cleanup execution
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			// Refresh sending deadline limit
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The central Hub explicitly closed this channel denoting removal
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Actually push physical characters across the wire seamlessly
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return // Wire broken, kill loop
			}

		case <-ticker.C:
			// Tick occurred, ping the client to prove we're still connected
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
