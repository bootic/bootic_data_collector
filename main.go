package main

import (
  "net/http"
  "datagram.io/udp"
  "datagram.io/udp/ws"
  "datagram.io/fanout"
  "log"
  "os"
)

func main() {
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
  log.Println("websocket server at " + ws_host + "/ws")
  
  // Push incoming UDP messages to multiple listeners ++++++++++++++++++
  // Push all events
  daemon.Subscribe(wshub.Notifier)
  
  fanoutObserver := fanout.NewZmq("tcp://127.0.0.1:6000")
  daemon.Subscribe(fanoutObserver.Notifier)
  log.Println("ZMQ fanout at tcp://127.0.0.1:6000")
  

  log.Fatal("HTTP server error: ", http.ListenAndServe(ws_host, nil))
}
