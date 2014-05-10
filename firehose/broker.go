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
  defunctClients chan chan []byte
  clients map[chan []byte]bool
}

func (broker *Broker) listen() {
  for {
    select {
    case s := <-broker.newClients:
      broker.clients[s] = true
      log.Printf("Client added. %d registered clients", len(broker.clients))
    case s := <-broker.defunctClients:
      // A client has dettached and we want to
      // stop sending them messages.
      delete(broker.clients, s)
      log.Printf("Removed client. %d registered clients", len(broker.clients))
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
      // log.Printf("Broadcast message to %d clients", len(broker.clients))
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
  rw.Header().Set("Access-Control-Allow-Origin", "*")

  messageChan := make(chan []byte)
  broker.newClients <- messageChan

  // Remove this client from the map of attached clients
  // when `EventHandler` exits.
  defer func() {
    fmt.Println("HERE.")
    broker.defunctClients <- messageChan
  }()

  // isten to connection close and un-register messageChan
  notify := rw.(http.CloseNotifier).CloseNotify()

  go func() {
      <-notify
      fmt.Println("HTTP connection just closed.")
      broker.defunctClients <- messageChan
  }()

  // block waiting or messages broadcast on this connection's messageChan
  for {
    // Write to the ResponseWriter
    fmt.Fprintf(rw, "data: %s\n\n", <-messageChan)
    f.Flush()
  }
}

func NewServer() (broker *Broker) {

  broker = &Broker{
    Notifier: make(data.EventsChannel, 1),
    newClients: make(chan chan []byte),
    defunctClients: make(chan chan []byte),
    clients: make(map[chan []byte]bool),
  }

  go broker.listen()

  return
}
