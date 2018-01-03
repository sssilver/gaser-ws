package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	pongPeriod = 60 * time.Second
	pingPeriod = 30 * time.Second
)

type Client struct {
	conn *websocket.Conn
	out  chan OutFrame
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
		out:  make(chan OutFrame),
	}
}

func (c *Client) tx() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case frame, ok := <-c.out: // Regular frames
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println("CloseMessage")
				return
			}

			// Send the outbound frame as JSON
			if err := c.conn.WriteJSON(frame); err != nil {
				log.Println("error: ", err)
				return
			}

		case <-ticker.C: // Ping frames
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) rx(gameInChannel chan InFrame) {
	defer func() {
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongPeriod))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(
			time.Now().Add(pongPeriod))
		return nil
	})

	for {
		frame := InFrame{}
		err := c.conn.ReadJSON(&frame)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				// The client has closed the connection
				log.Println("connection closed: ", err)
			} else {
				log.Println("error: ", err)
			}
			break
		}

		log.Println("got here")

		// Send the frame to the game
		gameInChannel <- frame
	}
}
