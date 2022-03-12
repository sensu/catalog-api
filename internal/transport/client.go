package transport

import (
	"time"

	"github.com/gorilla/websocket"
)

const (
	pingInterval = 30 * time.Second
	pongWait     = 60 * time.Second
	writeWait    = 10 * time.Second
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	transport *Transport
	conn      *websocket.Conn
	send      chan []byte
}

func NewClient(t *Transport, conn *websocket.Conn) Client {
	return Client{
		transport: t,
		conn:      conn,
		send:      make(chan []byte, 256),
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingInterval)
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

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
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
