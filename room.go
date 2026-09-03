package main

import (
	"fmt"
	"sync"
)

type room struct {
	mu      sync.RWMutex
	clients map[*client]struct{}
}

func (rm *room) add(c *client) {
	rm.mu.Lock()
	rm.clients[c] = struct{}{}
	rm.mu.Unlock()
}

func (rm *room) remove(c *client) {
	rm.mu.Lock()
	if _, ok := rm.clients[c]; ok {
		delete(rm.clients, c)
		close(c.send)
	}
	rm.mu.Unlock()
}

func (rm *room) broadcast(sender *client, msg []byte) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	var formatted []byte
	if sender == nil {
		// Server message; no sender
		formatted = fmt.Appendf(nil, "%v", string(msg))
	} else {
		formatted = fmt.Appendf(nil, "[%s] %v", sender.username, string(msg))
	}

	for client := range rm.clients {
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
			fmt.Printf("%s buffer full, dropping message\n", client.username)
		}
	}
}

func newRoom() *room {
	return &room{
		clients: make(map[*client]struct{}),
	}
}
