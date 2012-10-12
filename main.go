package main

import (
	"datagram.io/daemon"
	"datagram.io/ws"
  "net/http"
  "fmt"
  "log"
  "datagram.io/web"
  "datagram.io/db"
  "datagram.io/cmd"
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
		"setupdb":            db.SetupDB,
		"store-event":        cmd.StoreEvent,
		"daemon":             daemons,
		"help":               cmd.Help,
	}

	argc := len(os.Args)
	commandName := "help"

	if argc > 1 {
		commandName = os.Args[1]
	}

	var command func() error

	if command = commands[commandName]; command == nil {
		command = cmd.Help
	}

	if err := command(); err != nil {
		fmt.Println(err)
	}
}
