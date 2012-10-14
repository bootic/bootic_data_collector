package ws

import (
	"encoding/json"
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

func decodeEventIntoString(event *data.Event) (str string, err error) {
	bytes, err := json.Marshal(event)
	if err != nil {
		return
	}
	return string(bytes), err
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



func (this *Hub) Run() {
	for {
		select {
		case c := <-this.register:
			this.connections[c] = true
		case c := <-this.unregister:
			delete(this.connections, c)
			close(c.send)
		case m := <-this.broadcast:
			for c := range this.connections {
				select {
				case c.send <- m:
				default:
					delete(this.connections, c)
					close(c.send)
					go c.ws.Close()
				}
			}
		}
	}
}