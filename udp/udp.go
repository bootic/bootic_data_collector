package udp

import (
  "datagram.io/data"
  "github.com/bitly/go-simplejson"
  "log"
  "net"
)

type Daemon struct {
  Conn *net.UDPConn
  observers map[string][]data.EventsChannel
}

func NewDaemon(udpHostAndPort string) (daemon *Daemon, err error) {
  conn, err := createUDPListener(udpHostAndPort)

  if err != nil {
    return
  }
  
  daemon = &Daemon{
    Conn: conn,
    observers: make(map[string][]data.EventsChannel),
  }
  
  go daemon.ReceiveDatagrams()
  
  return
}

func (self *Daemon) Subscribe(observer data.EventsChannel) {
  self.observers["all"] = append(self.observers["all"], observer)
}

func (self *Daemon) SubscribeToType(observer data.EventsChannel, typeStr string) {
  self.observers[typeStr] = append(self.observers[typeStr], observer)
}

func (self *Daemon) Dispatch(event *simplejson.Json) {
  // Dispatch to global observers
  for _, observer := range self.observers["all"] {
    observer <- event
  }
  
  // Dispatch to type observers
  evtStr, _ := event.Get("type").String()
  for _, observer := range self.observers[evtStr] {
    observer <- event
  }
}

func (self *Daemon) ReceiveDatagrams() {

  for {
    buffer := make([]byte, 1024)
    
    if c, addr, err := self.Conn.ReadFromUDP(buffer); err != nil {

      log.Printf("blergh: %d byte datagram from %s with error %s\n", c, addr.String(), err.Error())
      return
    
    } else {
      
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

func createUDPListener(hostAndPort string) (conn *net.UDPConn, err error) {

	var udpaddr *net.UDPAddr
	if udpaddr, err = net.ResolveUDPAddr("udp4", hostAndPort); err != nil {
		return
	}

	conn, err = net.ListenUDP("udp4", udpaddr)

	return
}
