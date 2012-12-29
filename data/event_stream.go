package data

import (
	"encoding/json"
	"time"
	"fmt"
)

type EventStream struct {
	Events chan *Event
}

func jsonBytesIntoEvent(payload []byte) (event Event, err error) {
	err = json.Unmarshal(payload, &event)
	event.CreatedOn = time.Now()
	return
}

func (events *EventStream) WriteBytes(payload []byte) {
	go func() {

		event, err := jsonBytesIntoEvent(payload) //simplejson.NewJson([]byte(msg))

		if err != nil {
		  fmt.Println("Invalid JSON: " + string(payload))
		} else {
		  events.Events <- &event
		}
	}()
}

func NewEventStream() *EventStream {
	return &EventStream{
		Events: make(chan *Event),
	}
}
