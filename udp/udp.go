package udp

import (
  "datagram.io/data"
  "github.com/bitly/go-simplejson"
  "log"
  "net"
)

type Daemon struct {
  Conn *net.UDPConn
  Stream *data.EventStream
  observers []data.EventsChannel
}

func NewDaemon(udpHostAndPort string) (daemon *Daemon, err error) {
  conn, err := createUDPListener(udpHostAndPort)

  if err != nil {
    return
  }
  
  daemon = &Daemon{
    Conn: conn,
    Stream: data.NewEventStream(),
  }
  
  go daemon.ReceiveDatagrams()
  
  return
}

func (self *Daemon) Subscribe(observer data.EventsChannel) {
  self.observers = append(self.observers, observer)
}

func (self *Daemon) Dispatch(event *simplejson.Json) {
  for _, observer := range self.observers {
    observer <- event
  }
}

func (self *Daemon) ReceiveDatagrams() {

  for {
    buffer := make([]byte, 1024)
    
    if c, addr, err := self.Conn.ReadFromUDP(buffer); err != nil {
    
      log.Println("blergh: " + err.Error())
      return
    
    } else {
    
      log.Printf("received %d byte datagram from %s\n", c, addr.String())
      
      event, err := data.JsonBytesIntoEvent(buffer[:c])
      if err != nil {
        log.Println("Invalid JSON", err)
      } else {
        self.Dispatch(event)
      }
    }
  
  }

	panic("should never have got myself into this.")
}

func (self *Daemon) FilterByType(typeStr string) *data.EventStream {
  stream := data.NewEventStream()
  
  go func (s *data.EventStream) {
    for {
      event := <- self.Stream.Events
      eventType, _ := event.Get("type").String()
      
      if eventType == typeStr {
        stream.Events <- event
      }
    }
  }(stream)
  
  return stream
}

func createUDPListener(hostAndPort string) (conn *net.UDPConn, err error) {

	var udpaddr *net.UDPAddr
	if udpaddr, err = net.ResolveUDPAddr("udp4", hostAndPort); err != nil {
		return
	}

	conn, err = net.ListenUDP("udp4", udpaddr)

	return
}
