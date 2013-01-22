package data

import(
  "github.com/bitly/go-simplejson"
)

type EventsChannel chan *simplejson.Json

func JsonBytesIntoEvent(payload []byte) (event *simplejson.Json, err error) {
  event, err = simplejson.NewJson(payload)//json.Unmarshal(payload, &event)
  // event.CreatedOn = time.Now()
  return
}