package firehose

import (
  data "github.com/bootic/bootic_go_data"
  "log"
  "fmt"
  "net/http"
)

type Broker struct {
  Notifier data.EventsChannel
  newClients chan chan []byte
  clients map[chan []byte]bool
}

func (broker *Broker) listen() {
  for {
    select {
    case s := <-broker.newClients:
      broker.clients[s] = true
      log.Println("New client connected")
    case event := <-broker.Notifier:
      // Send event to all connected clients
      for clientMessageChan, _ := range broker.clients {
        json, err := data.EncodeJSON(event)
        if err != nil {
          log.Println("Error encoding event to JSON")
        } else {
          clientMessageChan <- json
        }
      }
      log.Printf("Broadcast message to %d clients", len(broker.clients))
    }
  }
  
}

func (broker *Broker) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
  // Make sure that the writer supports flushing.
  //
  f, ok := rw.(http.Flusher)

  if !ok {
    http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
    return
  }

  // Set the headers related to event streaming.
  rw.Header().Set("Content-Type", "text/event-stream")
  rw.Header().Set("Cache-Control", "no-cache")
  rw.Header().Set("Connection", "keep-alive")

  messageChan := make(chan []byte)
  broker.newClients <- messageChan

  for {
    // rw.Write(<-messageChan)
    // Write to the ResponseWriter, `w`.
    fmt.Fprintf(rw, "data: %s\n\n", <-messageChan)
    f.Flush()
  }
}

func NewServer() (broker *Broker) {

  broker = &Broker{
    Notifier: make(data.EventsChannel, 1),
    newClients: make(chan chan []byte),
    clients: make(map[chan []byte]bool),
  }

  go broker.listen()

  return
}
