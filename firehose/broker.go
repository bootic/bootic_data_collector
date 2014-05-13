package firehose

import (
	"fmt"
	data "github.com/bootic/bootic_go_data"
	"log"
	"net/http"
)

// A MessageChan is a channel of channels
// Each connection sends a channel of bytes to a global MessageChan
// The main broker listen() loop listens on new connections on MessageChan
// New event messages are broadcast to all registered connection channels
type MessageChan chan []byte

type Broker struct {
	Notifier       data.EventsChannel
	newClients     chan MessageChan
	defunctClients chan MessageChan
	clients        map[MessageChan]bool
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
			json, err := data.EncodeJSON(event)
			if err != nil {
				log.Println("Error encoding event to JSON")
			} else {
				for clientMessageChan, _ := range broker.clients {
					clientMessageChan <- json
				}
			}
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

	messageChan := make(MessageChan)
	broker.newClients <- messageChan

	// Remove this client from the map of attached clients
	// when `EventHandler` exits.
	defer func() {
		fmt.Println("HERE.")
		broker.defunctClients <- messageChan
	}()

	// "raw" query string option
	req.ParseForm()
	raw := len(req.Form["raw"]) > 0

	// Listen to connection close and un-register messageChan
	notify := rw.(http.CloseNotifier).CloseNotify()

	go func() {
		<-notify
		fmt.Println("HTTP connection just closed.")
		broker.defunctClients <- messageChan
	}()

	// block waiting or messages broadcast on this connection's messageChan
	for {
		// Write to the ResponseWriter
		if raw {
			// Raw JSON events, one per line
			fmt.Fprintf(rw, "%s\n", <-messageChan)
		} else {
			// Server Sent Events compatible
			fmt.Fprintf(rw, "data: %s\n\n", <-messageChan)
		}

		f.Flush()
	}
}

func NewServer() (broker *Broker) {

	broker = &Broker{
		Notifier:       make(data.EventsChannel, 1),
		newClients:     make(chan MessageChan),
		defunctClients: make(chan MessageChan),
		clients:        make(map[MessageChan]bool),
	}

	go broker.listen()

	return
}
