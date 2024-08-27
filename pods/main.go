package main

import (
	"context"
	"flag"
	"log"
	"noroom/pods/hub"
	"noroom/pods/server"
)

func main() {
	port := flag.Int("port", 6969, "port to use for listening")
	flag.Parse()

	hub, err := hub.NewHub(context.Background())
	if err != nil {
		log.Fatal("failed to create hub:", err)
	}

	srv := server.NewServer(hub)

	log.Println("listening on port", *port)
	if err := srv.Start(*port); err != nil {
		log.Fatal("server error:", err)
	}
}
