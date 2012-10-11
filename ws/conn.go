package ws

import (
	"code.google.com/p/go.net/websocket"
)

type Connection struct {
	// The websocket connection.
	ws *websocket.Conn

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
		Hub.Broadcast <- message
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

func WsHandler(ws *websocket.Conn) {
	c := &Connection{send: make(chan string, 256), ws: ws}
	Hub.register <- c
	defer func() { Hub.unregister <- c }()
	go c.writer()
	c.reader()
}
