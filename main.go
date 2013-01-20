package main

import (
  "datagram.io/udp"
  "datagram.io/udp/ws"
  // "datagram.io/db"
  "fmt"
  "log"
  "net/http"
  "os"
)

func daemons() (err error) {
  // Configure via env variables
  // WS and UDP hosts can be different, ex. UDP could be listening on a private IP while WS is public
	udp_host  := os.Getenv("DATAGRAM_IO_UDP_HOST")
	ws_host   := os.Getenv("DATAGRAM_IO_WS_HOST")
	
	// Start up UDP daemon +++++++++++++++++++++++++++++++++++++++++++++++
	daemon, err := udp.NewDaemon(udp_host)
	if err != nil {
	  panic(err)
	}

	// Setup Websockets hub ++++++++++++++++++++++++++++++++++++++++++++++
	wshub := ws.HandleWebsocketsHub("/ws")
	fmt.Println("websocket server at " + ws_host + "/ws")

	// Push incoming UDP messages to multiple listeners ++++++++++++++++++
	// Push all events
  wshub.Receive(daemon.Stream)
  // We can also filter events by type
	// wshub.Receive(daemon.FilterByType("pageview"))

	log.Fatal("HTTP server error: ", http.ListenAndServe(ws_host, nil))

	return nil
}

func main() {
  if err := daemons(); err != nil {
    fmt.Println(err)
  }
}
