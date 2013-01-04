package data

import (
  // "time"
	"fmt"
	"github.com/bitly/go-simplejson"
)

type EventStream struct {
	Events chan *simplejson.Json
}

func jsonBytesIntoEvent(payload []byte) (event *simplejson.Json, err error) {
	event, err = simplejson.NewJson(payload)//json.Unmarshal(payload, &event)
  // event.CreatedOn = time.Now()
	return
}

func (events *EventStream) WriteBytes(payload []byte) {
	go func() {

		event, err := jsonBytesIntoEvent(payload) //simplejson.NewJson([]byte(msg))

		if err != nil {
		  fmt.Println("Invalid JSON: " + string(payload))
		} else {
		  events.Events <- event
		}
	}()
}

func NewEventStream() *EventStream {
	return &EventStream{
		Events: make(chan *simplejson.Json),
	}
}
