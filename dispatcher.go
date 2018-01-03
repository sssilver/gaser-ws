package main

import (
	"io"
	"log"
	"net"
)

func game_tx(gameOutChannel chan InFrame, gameConn net.Conn) {
	for {
		log.Println("game_tx")
		// Take the outbound frame to the game
		frame := <-gameOutChannel

		// Forward it to the game server
		numWritten, err := gameConn.Write(frame.data) // TODO: Convert the entire frame to []byte and send it all
		// TODO: Or perhaps these frames should be a different type altogether
		if err != nil || numWritten < len(frame.data) { // See above TODO for len check
			log.Println("Error writing to the game: ", err)
			continue
		}
	}
}

func game_rx(gameInChannel chan OutFrame, gameConn net.Conn) {
	const (
		bufferSize = 4096
		chunkSize  = 1024
	)

	for {
		log.Println("game_rx")
		buffer := make([]byte, 0, bufferSize)
		chunk := make([]byte, chunkSize)

		for { // Read the TCP response by chunks
			numRead, err := gameConn.Read(chunk)

			if err != nil {
				if err != io.EOF {
					log.Println("Error reading from the game")
				}
				break
			}

			buffer = append(buffer, chunk[:numRead]...)
		}

		// Deserialize the buffer into a frame
		// TODO: These frames should prolly be different from the client ones
		gameInChannel <- OutFrame{data: buffer}
	}
}
