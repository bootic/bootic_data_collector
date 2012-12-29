package main

import (
	"datagram.io/daemon"
	"datagram.io/daemon/ws"
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
	udpEventStream := daemon.ReceiveDatagrams(udp_host)

	// newEvents := db.StoreEvents(udpEventStream)

	// Setup Websockets hub ++++++++++++++++++++++++++++++++++++++++++++++
	wshub := ws.HandleWebsocketsHub("/ws")
	fmt.Println("websocket server at " + ws_host + "/ws")

	// Push incoming UDP messages to multiple listeners ++++++++++++++++++
	wshub.Receive(udpEventStream)

	log.Fatal("HTTP server error: ", http.ListenAndServe(ws_host, nil))

	return nil
}

func main() {
  if err := daemons(); err != nil {
    fmt.Println(err)
  }
}
