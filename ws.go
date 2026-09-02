package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/coder/websocket"
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer c.CloseNow()

	ctx, cancel := context.WithTimeout(r.Context(), time.Minute)
	defer cancel()

	_, bytes, err := c.Read(ctx)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(bytes))
}
