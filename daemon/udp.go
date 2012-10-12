package daemon

import (
  "net"
  "fmt"
  "log"
)

func ReceiveDatagrams (hostAndPort string) EventStream {

	var conn *net.UDPConn
  
  conn, err := createUDPListener(hostAndPort)
  
  if err != nil {
  	panic("Could not create UDP listener")
  }

  fmt.Printf("Listener for UDP connections on %s\n", conn.LocalAddr().String())
  
  eventStream := newEventStream()
  
  go rcv(conn, eventStream)
  
  return *eventStream
  
}

func rcv (conn *net.UDPConn, eventStream *EventStream) {
	for {
	  buffer := make([]byte, 256)

  	if c, addr, err := conn.ReadFromUDP(buffer); err != nil {

  		fmt.Println("blergh: " + err.Error())
  		return

  	} else {
      
      log.Printf("received %d byte datagram from %s\n", c, addr.String())

      eventStream.writeBytes(buffer[:c])
  	}	
  	
	}
	panic("should never have got myself into this.")
}



func createUDPListener (hostAndPort string) (conn *net.UDPConn, err error) {

	var udpaddr *net.UDPAddr
	if udpaddr, err = net.ResolveUDPAddr("udp4", hostAndPort); err != nil {
		return
	}

  conn, err = net.ListenUDP("udp4", udpaddr)
 
	return
}

