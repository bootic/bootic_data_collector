package daemon


import(
	"encoding/json"
)

type EventStream struct {
  Events chan *db.Event
}

func jsonBytesIntoEvent(payload []byte) (err error, event db.Event) {
	err := json.Unmarshal(payload, &event)

	return
}

func (this *EventStream) writeBytes(payload []byte) {
  go func(){
    event, err := jsonBytesIntoEvent(payload)//simplejson.NewJson([]byte(msg))
    if err != nil {
      panic("Invalid JSON: " + payload)
    }
   	
    this.Events <- event
  }()
}

func newEventStream() *EventStream {
  return &EventStream{
    Events: make(chan *db.Event),
  }
}