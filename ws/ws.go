package ws

import (
	"code.google.com/p/go.net/websocket"
	"net/http"
)

type Connection struct {
	// The websocket connection.
	ws *websocket.Conn
  hub *Hub
  
	// Buffered channel of outbound messages.
	send chan string
}

func (c *Connection) reader() {
	for {
		var message string
		err := websocket.Message.Receive(c.ws, &message)
		if err != nil {
			break
		}
		c.hub.broadcast <- message
	}
	c.ws.Close()
}

func (c *Connection) writer() {
	for message := range c.send {
		err := websocket.Message.Send(c.ws, message)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}

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

func (h *Hub) Receive(input chan string) {
  go func() {
    for {
     h.broadcast <- <- input
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

func HandleWebsocketsHub (path string) *Hub {

  hub := NewHub()
  
  http.Handle(path, websocket.Handler(func(ws *websocket.Conn) {
    c := &Connection{send: make(chan string, 256), ws: ws, hub: hub}
  	hub.register <- c
  	defer func() { hub.unregister <- c }()
  	go c.writer()
  	c.reader()
  }))
  
  return hub
}
