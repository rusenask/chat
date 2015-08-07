package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

//client representsa single chatting user
type client struct {
	// socket is the web socket for this client
	socket *websocket.Conn
	// send is a channel on which messages are sent
	send chan []byte
	// room is the room this client is chatting in
	room *room
}

func (c *client) read() {
	fmt.Println("reading")
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			fmt.Println(msg)
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	fmt.Println("writing")
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
