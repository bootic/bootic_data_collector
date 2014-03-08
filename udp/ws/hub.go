package ws

import (
  data "github.com/bootic/bootic_go_data"
)

type Hub struct {
  // Registered connections.
  connections map[*Connection]bool

  // Inbound messages from the connections.
  Notifier data.EventsChannel

  // Register requests from the connections.
  register chan *Connection

  // Unregister requests from connections.
  unregister chan *Connection
}

func NewHub() *Hub {
  h := &Hub{
    Notifier:    make(data.EventsChannel),
    register:    make(chan *Connection),
    unregister:  make(chan *Connection),
    connections: make(map[*Connection]bool),
  }

  go h.Run()

  return h
}

func (this *Hub) Run() {
  for {
    select {
    case c := <-this.register:
      this.connections[c] = true
    case c := <-this.unregister:
      delete(this.connections, c)
      close(c.send)
    case event := <-this.Notifier:
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
