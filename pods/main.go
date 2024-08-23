package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()

	central, err := NewCentral(context.Background())
	if err != nil {
		log.Fatal("NewCentral:", err)
	}
	defer central.Close()

	go central.Start()

	http.HandleFunc("/", serveTest)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(central, w, r)
	})

	fmt.Println("listening on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveWs(central *Central, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(central, conn)
	go client.writePump()
	go client.readPump()
}

func serveTest(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path == "/" {
		http.ServeFile(w, r, "test.html")
		return
	}

	if r.URL.Path == "/test.js" {
		http.ServeFile(w, r, "test.js")
		return
	}

	http.Error(w, "Not found", http.StatusNotFound)
	return
}
