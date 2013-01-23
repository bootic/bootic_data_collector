package main

import (
  "net/http"
  "datagram.io/udp"
  "datagram.io/udp/ws"
  // "datagram.io/db"
  "datagram.io/redis_stats"
  "log"
  "os"
)

func daemons() (err error) {
  // Configure via env variables
  // WS and UDP hosts can be different, ex. UDP could be listening on a private IP while WS is public
  udp_host  := os.Getenv("DATAGRAM_IO_UDP_HOST")
  ws_host   := os.Getenv("DATAGRAM_IO_WS_HOST")
  redis_host:= os.Getenv("REDIS_HOST")
  
  // Start up UDP daemon +++++++++++++++++++++++++++++++++++++++++++++++
  daemon, err := udp.NewDaemon(udp_host)
  if err != nil {
    panic(err)
  }
  
  // Track some event types to redis +++++++++++++++++++++++++++++++++++
  if redis_host != "" {
    tracker, err := redis_stats.NewTracker(redis_host)

    if err != nil {
      panic(err)
    }
    
    log.Println("Using redis", redis_host, tracker)
    // Track pageview events as redis increment time series
    daemon.SubscribeToType(tracker.Notifier, "pageview")
  }
  
	// Setup Websockets hub ++++++++++++++++++++++++++++++++++++++++++++++
	wshub := ws.HandleWebsocketsHub("/ws")
	log.Println("websocket server at " + ws_host + "/ws")

	// Push incoming UDP messages to multiple listeners ++++++++++++++++++
	// Push all events
  daemon.Subscribe(wshub.Notifier)
  // We can also filter events by type

	log.Fatal("HTTP server error: ", http.ListenAndServe(ws_host, nil))

	return nil
}

func main() {
  if err := daemons(); err != nil {
    log.Println(err)
  }
}
