package daemon

import (
  "net"
  "fmt"
)

func ReceiveDatagrams (hostAndPort string) {

	var conn *net.UDPConn
  
  conn, err := createUDPListener(hostAndPort)
  
  if err != nil {
  	panic("Could not create UDP listener")
  }

  fmt.Printf("Listener for UDP connections on %s\n", conn.LocalAddr().String())

  go rcv(conn)
}

func rcv (conn *net.UDPConn) {
	for {
	  buffer := make([]byte, 256)

  	if c, addr, err := conn.ReadFromUDP(buffer); err != nil {

  		fmt.Println("blergh: " + err.Error())
  		return

  	} else {

  		fmt.Printf("%d byte datagram received from %s\n\n", c, addr.String())
  		fmt.Printf("\t\"%s\"\n\n", string(buffer[:c]))
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

