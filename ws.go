package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/coder/websocket"
)

var clientCount = 0

type client struct {
	c  *websocket.Conn
	id int
}

func newClient(c *websocket.Conn, id int) *client {
	return &client{c, id}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer c.CloseNow()

	cl := newClient(c, clientCount)
	clientCount++

	fmt.Printf("Client [%d] connected\n", cl.id)

	ctx := context.Background()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, bytes, err := c.Read(ctx)
			if err != nil {
				// Normal disconnect
				if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
					websocket.CloseStatus(err) == websocket.StatusGoingAway ||
					errors.Is(err, io.EOF) {
					fmt.Printf("Client [%d] disconnected\n", cl.id)
					return
				}

				fmt.Println(err)
			}

			fmt.Printf("[%d] %v\n", cl.id, string(bytes))
		}
	}
}
