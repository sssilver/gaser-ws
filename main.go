package main

import (
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	readBufferSize  = 1024
	writeBufferSize = 1024
	maxMessageSize  = 1024
)

func handler(clients map[*Client]bool, gameInChannel chan InFrame, w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  readBufferSize,
		WriteBufferSize: writeBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade: ", err)
	}

	client := NewClient(conn)
	clients[client] = true

	go client.tx()
	go client.rx(gameInChannel)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Set of all connected clients
	clients := make(map[*Client]bool)

	// Create the game inbound channel
	gameInChannel := make(chan InFrame)
	gameOutChannel := make(chan OutFrame)

	// Connect to the game server
	// TODO: This is temporary; in reality this should be some protobuf over TCP
	gameConn, _ := net.Pipe()

	// Create the game dispatchers
	go game_tx(gameOutChannel, gameConn)
	go game_rx(gameInChannel, gameConn)

	// Set up the listener
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(clients, gameInChannel, w, r)
	})

	// Start listening
	// TODO: Read the listen host from command-line or config or something
	err := http.ListenAndServe("localhost:4000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
