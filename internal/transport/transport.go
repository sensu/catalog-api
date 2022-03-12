package transport

import (
	"context"
)

type Transport struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewTransport() Transport {
	return Transport{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (t *Transport) Broadcast(msg []byte) {
	for client, connected := range t.clients {
		if connected {
			client.send <- msg
		}
	}
}

func (t *Transport) Register(client *Client) {
	t.register <- client
}

func (t *Transport) Start(ctx context.Context) {
	for {
		select {
		case client := <-t.register:
			t.clients[client] = true
		case client := <-t.unregister:
			if _, ok := t.clients[client]; ok {
				delete(t.clients, client)
				close(client.send)
			}
		}
	}
}
