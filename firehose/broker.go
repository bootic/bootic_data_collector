package firehose

import (
	"fmt"
	data "github.com/bootic/bootic_go_data"
	"log"
	"net/http"
)

// The Firehose server implements a Server Sent Events endpoint.
// It subscribes to the global event stream and exposes them as a long-lived, authenticated
// HTTP response.
// Example usage:
//
//     $ curl -u user:pass server.com
//     data: {"type": "pageview", "time": 123456789, "data": {...}}
//     data: {"type": "order", "time": 123456799, "data": {...}}

// A MessageChan is a channel of channels
// Each connection sends a channel of bytes to a global MessageChan
// The main broker listen() loop listens on new connections on MessageChan
// New event messages are broadcast to all registered connection channels
type MessageChan chan []byte

// A Broker holds open client connections,
// listens for incoming events on its Notifier channel
// and broadcast event data to all registered connections
type Broker struct {

	// Events are pushed to this channel by the main UDP daemon
	Notifier data.EventsChannel

	// New client connections
	newClients chan MessageChan

	// Closed client connections
	closingClients chan MessageChan

	// Client connections registry
	clients map[MessageChan]bool
}

// Listen on different channels and act accordingly
func (broker *Broker) listen() {
	for {
		select {
		case s := <-broker.newClients:

			// A new client has connected.
			// Register their message channel
			broker.clients[s] = true
			log.Printf("Client added. %d registered clients", len(broker.clients))
		case s := <-broker.closingClients:

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

// Implement the http.Handler interface.
// This allows us to wrap HTTP handlers (see auth_handler.go)
// http://golang.org/pkg/net/http/#Handler
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

	// Each connection registers its own message channel with the Broker's connections registry
	messageChan := make(MessageChan)

	// Signal the broker that we have a new connection
	broker.newClients <- messageChan

	// Remove this client from the map of attached clients
	// when `EventHandler` exits.
	defer func() {
		fmt.Println("HERE.")
		broker.closingClients <- messageChan
	}()

	// "raw" query string option
	// If provided, send raw JSON lines instead of SSE-compliant strings.
	req.ParseForm()
	raw := len(req.Form["raw"]) > 0

	// Listen to connection close and un-register messageChan
	notify := rw.(http.CloseNotifier).CloseNotify()

	go func() {
		<-notify
		fmt.Println("HTTP connection just closed.")
		broker.closingClients <- messageChan
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

// Server factory
func NewServer() (broker *Broker) {

	broker = &Broker{
		Notifier:       make(data.EventsChannel, 1),
		newClients:     make(chan MessageChan),
		closingClients: make(chan MessageChan),
		clients:        make(map[MessageChan]bool),
	}

	go broker.listen()

	return
}
