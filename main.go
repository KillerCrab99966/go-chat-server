package main

import (
	_ "embed"
	"log"
	"net/http"
	"strings"
	"time"
)

//go:embed short-wordlist.txt
var wordlist string

func main() {
	gen := newCodeGenerator(strings.Split(wordlist, "\n"), 4)

	rm, err := newRoom("Room1", gen)
	if err != nil {
		log.Fatal(err)
	}

	_, m := newApp(6, time.Second*30, rm)

	startServer(m)
}

func startServer(mux *http.ServeMux) {
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Server listening on", server.Addr)
	log.Fatal(server.ListenAndServe())
}
