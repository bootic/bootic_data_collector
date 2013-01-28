package fanout

import (
  "datagram.io/data"
  zmq "github.com/alecthomas/gozmq"
  "bytes"
)

type ZMQObserver struct {
  socket zmq.Socket
  Notifier data.EventsChannel
}

func (observer *ZMQObserver) listen () {
  for {
    event := <- observer.Notifier
    evtBytes, err := data.Encode(event)
    if err != nil {
      break
    }
    evtType, _ := event.Get("type").String()
    
    observer.dispatch(evtType, evtBytes)
  }
}

func (observer *ZMQObserver) dispatch(topicStr string, evtBytes []byte) {
  
  topic := []byte(topicStr)
  
  a := [][]byte{ topic, evtBytes } 
  
  topic_and_message := bytes.Join(a, []byte{' '})
  
  observer.socket.Send(topic_and_message, 0)
}

func NewZmq(host string) (observer *ZMQObserver) {
  context, _ := zmq.NewContext()
  socket, _ := context.NewSocket(zmq.PUB)
  
  socket.Bind(host)
  
  observer = &ZMQObserver{
    socket: socket,
    Notifier: make(data.EventsChannel, 1),
  }
  
  go observer.listen()
  
  return
}
