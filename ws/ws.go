package ws

import (
	"code.google.com/p/go.net/websocket"
	"net/http"
	"fmt"
  "strings"
)

type Connection struct {
	// The websocket connection.
	ws *websocket.Conn
  hub *Hub
  
	// Buffered channel of outbound messages.
	send chan string
	
	// Filters
	tags []string
}

func (c *Connection) reader() {
  tagsQuery := c.ws.Request().URL.Query().Get("tags")
  var tags []string
  
  if tagsQuery != ""{
    tags = strings.Split(tagsQuery, ",")
    c.tags = append(c.tags, tags...)
  }
  
  fmt.Println("ws [conn] initialized with", c.tags)
  
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
		fmt.Println("WS send " + message)
		if err != nil {
			break
		}
	}
	c.ws.Close()
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
