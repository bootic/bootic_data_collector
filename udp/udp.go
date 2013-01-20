package udp

import (
	"datagram.io/data"
	"log"
	"net"
)

type Daemon struct {
  Conn *net.UDPConn
  Stream *data.EventStream
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

func (self *Daemon) ReceiveDatagrams() {

  for {
    buffer := make([]byte, 1024)
    
    if c, addr, err := self.Conn.ReadFromUDP(buffer); err != nil {
    
      log.Println("blergh: " + err.Error())
      return
    
    } else {
    
      log.Printf("received %d byte datagram from %s\n", c, addr.String())
      
      self.Stream.WriteBytes(buffer[:c])
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
