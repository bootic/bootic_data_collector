package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"datagram.io/cmd"
	"datagram.io/db"
	datagramhttp "datagram.io/http"
)

func httpServer() (err error) {

	router := datagramhttp.Router()
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":5555", nil))

	return nil
}

func help() error {
	fmt.Println("commands available: http, generate-migration, help(this)")
	return nil
}

func storeEvent() (err error) {

	// if there is not specified event, return an error
	// read the event from the command line
	// save the event
	// return any error

	if len(os.Args) != 3 {
		err = fmt.Errorf("no event description specified")
		return
	}

	desc := os.Args[2]
	tags := []string{"tag1", "tag2", "tag3"}
	event := db.Event{Desc: desc, Tags: tags}

	err = db.StoreEvent(&event)

	return
}

func main() {

	db.Init()

	commands := map[string]func() error{
		"generate-migration": cmd.GenerateMigration,
		"setupdb":            db.SetupDB,
		"store-event":        storeEvent,
		"http":               httpServer,
		"help":               help,
	}

	argc := len(os.Args)
	commandName := "help"

	if argc > 1 {
		commandName = os.Args[1]
	}

	var command func() error

	if command = commands[commandName]; command == nil {
		command = help
	}

	if err := command(); err != nil {
		fmt.Println(err)
	}
}
