package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
)

var clientCount = 0

type client struct {
	conn *websocket.Conn
	id   int
	send chan []byte
}

func newClient(c *websocket.Conn, id int) *client {
	return &client{
		conn: c,
		id:   id,
		send: make(chan []byte, 256),
	}
}

type hub struct {
	mu      sync.RWMutex
	clients map[*client]struct{}
}

func newHub() *hub {
	return &hub{
		clients: make(map[*client]struct{}),
	}
}

func (h *hub) add(c *client) {
	h.mu.Lock()
	h.clients[c] = struct{}{}
	h.mu.Unlock()
}

func (h *hub) remove(c *client) {
	h.mu.Lock()
	if _, ok := h.clients[c]; ok {
		delete(h.clients, c)
		close(c.send)
	}
	h.mu.Unlock()
}

func (h *hub) broadcast(sender *client, msg []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var formatted []byte
	if sender == nil {
		// Server message; no sender
		formatted = fmt.Appendf(nil, "[] %v", string(msg))
	} else {
		formatted = fmt.Appendf(nil, "[%d] %v", sender.id, string(msg))
	}

	for client := range h.clients {
		// Do not send message to sender
		if sender != nil {
			if client.id == sender.id {
				continue
			}
		}

		select {
		case client.send <- formatted:
		default:
			// Buffer full
			fmt.Printf("Client %d buffer full, dropping message\n", client.id)
		}
	}
}

func (h *hub) wsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer c.CloseNow()

	cl := newClient(c, clientCount)
	clientCount++

	h.add(cl)
	defer h.remove(cl)

	go monitorMsgs(cl)

	msg := fmt.Sprintf("Client %d connected", cl.id)
	fmt.Printf("%s\n", msg)
	h.broadcast(nil, []byte(msg))

	ctx := r.Context()
	for {
		_, msg, err := c.Read(ctx)
		if err != nil {
			// Normal disconnect
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
				websocket.CloseStatus(err) == websocket.StatusGoingAway ||
				errors.Is(err, io.EOF) {

				msg := fmt.Sprintf("Client %d disconnected", cl.id)
				fmt.Printf("%s\n", msg)
				h.broadcast(nil, []byte(msg))

				return
			}

			fmt.Println(err)
		}

		fmt.Printf("[%d] %v\n", cl.id, string(msg))
		h.broadcast(cl, msg)
	}
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
