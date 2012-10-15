package main

import (
	"datagram.io/cmd"
	"datagram.io/daemon"
	"datagram.io/daemon/web"
	"datagram.io/daemon/ws"
	"datagram.io/db"
	"fmt"
	"log"
	"net/http"
	"os"
)

const hostAndPort = "localhost:5555"

func daemons() (err error) {

	// Start up UDP daemon +++++++++++++++++++++++++++++++++++++++++++++++
	udpEventStream := daemon.ReceiveDatagrams(hostAndPort)

	newEvents := db.StoreEvents(udpEventStream)

	// Setup Websockets hub ++++++++++++++++++++++++++++++++++++++++++++++
	wshub := ws.HandleWebsocketsHub("/ws")
	fmt.Println("websocket server at " + hostAndPort + "/ws")

	// Push incoming UDP messages to multiple listeners ++++++++++++++++++
	wshub.Receive(newEvents)

	router := web.Router()
	http.Handle("/", router)

	fmt.Println("serving HTTP at " + hostAndPort + "/")
	log.Fatal("HTTP server error: ", http.ListenAndServe(hostAndPort, nil))

	return nil
}

func main() {

	db.Init()

	commands := map[string]func() error{
		"setupdb":     db.SetupDB,
		"store-event": cmd.StoreEvent,
		"daemons":     daemons,
		"help":        cmd.ExplicitCallForHelp,
	}

	argc := len(os.Args)
	commandName := "help"

	if argc > 1 {
		commandName = os.Args[1]
	}

	var command func() error

	if command = commands[commandName]; command == nil {
		command = cmd.MissingCommandHelp
	}

	if err := command(); err != nil {
		fmt.Println(err)
	}
}
