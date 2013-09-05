package main

import (
  "net/http"
  "datagram.io/udp"
  "datagram.io/udp/ws"
  "datagram.io/fanout"
  "log"
  "flag"
)

func main() {
  var(
    udpHost string
    wsHost string
    zmqAddress string
  )
  
  // WS and UDP hosts can be different, ex. UDP could be listening on a private IP while WS is public
  flag.StringVar(&zmqAddress, "zmqsocket", "tcp://127.0.0.1:6000", "ZMQ socket address to send events to")
  flag.StringVar(&wsHost, "wshost", "localhost:5555", "Websocket host:port")
  flag.StringVar(&udpHost, "udphost", "localhost:5555", "host:port to bind for UDP datagrams")
  
  flag.Parse()
  
  // Start up UDP daemon +++++++++++++++++++++++++++++++++++++++++++++++
  daemon, err := udp.NewDaemon(udpHost)
  if err != nil {
    panic(err)
  }
  
  // Setup Websockets hub ++++++++++++++++++++++++++++++++++++++++++++++
  wshub := ws.HandleWebsocketsHub("/ws")
  log.Println("websocket server at " + wsHost + "/ws")
  
  // Push incoming UDP messages to multiple listeners ++++++++++++++++++
  // Push all events
  daemon.Subscribe(wshub.Notifier)
  
  // Start up PUB ZMQ client
  fanoutObserver := fanout.NewZmq(zmqAddress)
  
  // Push incoming UDP events down ZMQ pub/sub socket
  daemon.Subscribe(fanoutObserver.Notifier)
  log.Println("ZMQ fanout at", zmqAddress)
  

  log.Fatal("HTTP server error: ", http.ListenAndServe(wsHost, nil))
}
