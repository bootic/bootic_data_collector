package ws

import (
	"datagram.io/data"
)

type Hub struct {
	// Registered connections.
	connections map[*Connection]bool

	// Inbound messages from the connections.
	broadcast chan *data.Event

	// Register requests from the connections.
	register chan *Connection

	// Unregister requests from connections.
	unregister chan *Connection
}

func NewHub () (*Hub) {
 h := &Hub{
 	broadcast:   make(chan *data.Event),
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
      h.broadcast <- event
    }
  }()
}

func (this *Hub) Run() {
	for {
		select {
		case c := <-this.register:
			this.connections[c] = true
		case c := <-this.unregister:
			delete(this.connections, c)
			close(c.send)
		case event := <-this.broadcast:
			for c := range this.connections {
				select {
				case c.send <- event:
				default:
					delete(this.connections, c)
					close(c.send)
					go c.ws.Close()
				}
			}
		}
	}
}