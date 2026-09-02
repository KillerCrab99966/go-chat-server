package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/coder/websocket"
)

type client struct {
	conn     *websocket.Conn
	id       int
	username string
	send     chan []byte
}

func newClient(c *websocket.Conn, id int, username string) *client {
	return &client{
		conn:     c,
		id:       id,
		username: username,
		send:     make(chan []byte, 256),
	}
}

func randomiseUsername(raw string) string {
	id := rand.IntN(8999) + 1000
	return fmt.Sprintf("%s%d", raw, id)
}
