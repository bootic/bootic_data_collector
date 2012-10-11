package daemon

import (
  "net"
	"fmt"
)

func ReceiveDatagrams () (err error) {

	var conn *net.UDPConn

  if conn, err = createUDPListener("192.168.100.245:5555"); err != nil {
  	return err
  }

  fmt.Printf("Listener for UDP connections on %s\n", conn.LocalAddr().String())

  for {  	
  	rcv(conn)
  }

  panic("should never have got myself into this.")
}

func rcv (conn *net.UDPConn) {
		
	buffer := make([]byte, 256)

	if c, addr, err := conn.ReadFromUDP(buffer); err != nil {
		
		fmt.Println("blergh: " + err.Error())
		return
	
	} else {
		
		fmt.Printf("%d byte datagram received from %s\n\n", c, addr.String())
		fmt.Printf("\t\"%s\"\n\n", string(buffer[:c]))
	}	
}

func createUDPListener (hostAndPort string) (conn *net.UDPConn, err error) {

	var udpaddr *net.UDPAddr
	if udpaddr, err = net.ResolveUDPAddr("udp4", hostAndPort); err != nil {
		return
	}

  conn, err = net.ListenUDP("udp4", udpaddr)
 
	return
}

