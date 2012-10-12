package ws

import (
	"code.google.com/p/go.net/websocket"
	"net/http"
	"encoding/json"
	"fmt"
	"datagram.io/data"
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
		fmt.Println("AAAAA " + message)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}

func decodeEventIntoString(event *data.Event) (str string, err error) {
	bytes, err := json.Marshal(event)
	if err != nil {
		return
	}
	return string(bytes), err
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
