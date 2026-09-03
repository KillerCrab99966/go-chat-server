package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

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

func monitorMsgs(cl *client) {
	for msg := range cl.send {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		err := cl.conn.Write(ctx, websocket.MessageText, msg)

		cancel()

		if err != nil {
			cl.conn.Close(websocket.StatusGoingAway, "write failed")
			return
		}
	}
}
