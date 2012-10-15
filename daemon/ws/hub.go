package ws

import (
	"datagram.io/data"
)

type WebsocketHub struct {
	// Registered connections.
	connections map[*WebsocketConnection]bool

	// Inbound messages from the connections.
	broadcast chan string

	// Register requests from the connections.
	register chan *WebsocketConnection

	// Unregister requests from connections.
	unregister chan *WebsocketConnection
}

func NewWebsocketHub () (*WebsocketHub) {
 h := &WebsocketHub{
 	broadcast:   make(chan string),
 	register:    make(chan *WebsocketConnection),
 	unregister:  make(chan *WebsocketConnection),
 	connections: make(map[*WebsocketConnection]bool),
 }
 
 go h.Run()
 
 return h
}

func (h *WebsocketHub) Receive(eventStream *data.EventStream) {
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