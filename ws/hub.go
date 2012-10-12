package ws

import (
	"datagram.io/data"
)

type Hub struct {
	// Registered connections.
	connections map[*Connection]bool

	// Inbound messages from the connections.
	broadcast chan string

	// Register requests from the connections.
	register chan *Connection

	// Unregister requests from connections.
	unregister chan *Connection
}

func NewHub () (*Hub) {
 h := &Hub{
 	broadcast:   make(chan string),
 	register:    make(chan *Connection),
 	unregister:  make(chan *Connection),
 	connections: make(map[*Connection]bool),
 }
 
 go h.Run()
 
 return h
}

func (h *Hub) Receive(eventStream *data.EventStream) {
  go func() {
    for {
      event := <- eventStream.Events
      msg, err := decodeEventIntoString(event)
      if(err != nil) {
        break;
      }
      h.broadcast <- msg
    }
  }()
}