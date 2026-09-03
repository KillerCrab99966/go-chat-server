package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type app struct {
	// Key is room.code.
	rooms map[string]*room

	// Key is token.value.
	tokens   map[string]token
	tokenLen int
	tokenTTL time.Duration
}

func (a *app) wsHandler(w http.ResponseWriter, r *http.Request) {
	// Validate token (exists & not expired)
	tknParam := r.URL.Query().Get("token")
	tkn, ok := a.tokens[tknParam]
	if !ok || tkn.created.Add(a.tokenTTL).Before(time.Now()) {
		// Delete expired token
		delete(a.tokens, tknParam)

		http.Error(w, "Token invalid", http.StatusUnauthorized)
		return
	}

	// Delete token and forward request
	rm := tkn.room
	delete(a.tokens, tkn.value)
	rm.wsHandler(w, r)
}

func (a *app) roomJoinHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get code and associated room
	code := r.FormValue("code")
	if code == "" {
		http.Error(w, "Room code required", http.StatusBadRequest)
		return
	}

	rm, ok := a.rooms[code]
	if !ok {
		http.Error(w, "Room code invalid", http.StatusUnauthorized)
		return
	}

	// Create and store token
	tkn, err := newToken(rm, a.tokenLen)
	if err != nil {
		log.Printf("ERROR: crypto/rand failed to generate token: %v", err)

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	a.tokens[tkn.value] = tkn
	log.Printf("Token generated: %s\n", tkn.value)

	res := roomJoinResponse{
		Room:  rm.name,
		Token: tkn.value,
	}
	json.NewEncoder(w).Encode(res)
}

type roomJoinResponse struct {
	// The name of the Room.
	Room string

	// The generated Token/OTP.
	Token string
}

func newApp(tokenLen int, tokenTTL time.Duration, rooms ...*room) (*app, *http.ServeMux) {
	// Assign the rooms' codes as keys
	roomsMap := make(map[string]*room)
	for _, room := range rooms {
		roomsMap[room.code] = room
		fmt.Printf("%s: %s\n", room.name, room.code)
	}

	a := &app{
		rooms:    roomsMap,
		tokens:   make(map[string]token),
		tokenLen: tokenLen,
		tokenTTL: tokenTTL,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", a.wsHandler)
	mux.HandleFunc("POST /api/room/join", a.roomJoinHandler)

	return a, mux
}
