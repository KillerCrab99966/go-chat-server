package main

import (
	"fmt"
	"log"
	"sync"
)

type room struct {
	mu      sync.RWMutex
	name    string
	clients map[*client]struct{}
	code    string
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
			log.Printf("%s buffer full, dropping message\n", client.username)
		}
	}
}

func newRoom(name string, gen *codeGenerator) (*room, error) {
	code, err := gen.generateCode()
	if err != nil {
		return nil, fmt.Errorf("generating code: %w", err)
	}

	return &room{
		name:    name,
		clients: make(map[*client]struct{}),
		code:    code,
	}, nil
}
