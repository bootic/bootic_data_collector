package cmd

import (
	"os"
	"datagram.io/data"
	"datagram.io/db"
	"fmt"
)

func StoreEvent() (err error) {

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
	event := data.Event{Desc: desc, Tags: tags}

	err = db.StoreEvent(&event)

	return
}