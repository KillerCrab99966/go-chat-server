package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

var clientCount = 0

func (rm *room) wsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer c.CloseNow()

	// Avoid username conflicts
	rawUsername := r.URL.Query().Get("username")
	if rawUsername == "" {
		rawUsername = "User"
	}

	cl := newClient(c, clientCount, randomiseUsername(rawUsername))
	clientCount++

	rm.add(cl)
	defer rm.remove(cl)

	go monitorMsgs(cl)

	msg := fmt.Sprintf("%s connected to %s", cl.username, rm.name)
	log.Printf("%s\n", msg)
	rm.broadcast(nil, []byte(msg))

	ctx := r.Context()
	for {
		_, msg, err := c.Read(ctx)
		if err != nil {
			// Normal disconnect
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
				websocket.CloseStatus(err) == websocket.StatusGoingAway ||
				errors.Is(err, io.EOF) {

				msg := fmt.Sprintf("%s disconnected from %s", cl.username, rm.name)
				log.Printf("%s\n", msg)
				rm.broadcast(nil, []byte(msg))

				return
			}

			fmt.Println(err)
		}

		log.Printf("[%s] %v\n", cl.username, string(msg))
		rm.broadcast(cl, msg)
	}
}
