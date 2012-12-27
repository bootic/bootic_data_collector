package main

import (
	"datagram.io/daemon"
	"datagram.io/daemon/ws"
  // "datagram.io/db"
	"fmt"
	"log"
	"net/http"
)

const hostAndPort = "localhost:5555"

func daemons() (err error) {

	// Start up UDP daemon +++++++++++++++++++++++++++++++++++++++++++++++
	udpEventStream := daemon.ReceiveDatagrams(hostAndPort)

	// newEvents := db.StoreEvents(udpEventStream)

	// Setup Websockets hub ++++++++++++++++++++++++++++++++++++++++++++++
	wshub := ws.HandleWebsocketsHub("/ws")
	fmt.Println("websocket server at " + hostAndPort + "/ws")

	// Push incoming UDP messages to multiple listeners ++++++++++++++++++
	wshub.Receive(udpEventStream)

  // router := web.Router()
  // http.Handle("/", router)
	log.Fatal("HTTP server error: ", http.ListenAndServe(hostAndPort, nil))

	return nil
}

func main() {
  if err := daemons(); err != nil {
    fmt.Println(err)
  }
}
