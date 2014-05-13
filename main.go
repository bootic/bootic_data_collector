package main

import (
	"bootic_data_collector/fanout"
	"bootic_data_collector/firehose"
	"bootic_data_collector/udp"
	"bootic_data_collector/udp/ws"
	"flag"
	"log"
	"net/http"
)

func main() {
	var (
		udpHost           string
		wsHost            string
		zmqAddress        string
		sseHost           string
		globalAccessToken string
	)

	// WS and UDP hosts can be different, ex. UDP could be listening on a private IP while WS is public
	flag.StringVar(&zmqAddress, "zmqsocket", "tcp://127.0.0.1:6000", "ZMQ socket address to send events to")
	flag.StringVar(&wsHost, "wshost", "localhost:5555", "Websocket host:port")
	flag.StringVar(&udpHost, "udphost", "localhost:5555", "host:port to bind for UDP datagrams")
	flag.StringVar(&sseHost, "ssehost", "localhost:5556", "host:port to bind for Server Sent Events")
	flag.StringVar(&globalAccessToken, "accesstoken", "", "Access token required to connect to SSE events")

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
	zmqObserver := fanout.NewZmq(zmqAddress)

	// Push incoming UDP events down ZMQ pub/sub socket
	daemon.Subscribe(zmqObserver.Notifier)

	// Firehose subscriber
	firehoseDaemon := firehose.NewServer()

	daemon.Subscribe(firehoseDaemon.Notifier)

	log.Println("ZMQ fanout at", zmqAddress)

	authenticatedFirehose := firehose.NewAuthHandler(firehoseDaemon, globalAccessToken)

	go http.ListenAndServe(sseHost, authenticatedFirehose)

	log.Println("Server Sent Events at", sseHost)
	log.Fatal("HTTP server error: ", http.ListenAndServe(wsHost, nil))
}
