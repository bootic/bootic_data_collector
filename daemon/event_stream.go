package daemon


import(
	"encoding/json"
  "datagram.io/db"
)

type EventStream struct {
  Events chan *db.Event
}

func jsonBytesIntoEvent(payload []byte) (event db.Event, err error) {
	err = json.Unmarshal(payload, &event)
	return
}

func (events *EventStream) writeBytes(payload []byte) {
  go func(){
    
    event, err := jsonBytesIntoEvent(payload)//simplejson.NewJson([]byte(msg))
    
    if err != nil {
      panic("Invalid JSON: " + string(payload))
    }
   	
    events.Events <- &event
  }()
}

func newEventStream() *EventStream {
  return &EventStream{
    Events: make(chan *db.Event),
  }
}