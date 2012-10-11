package daemon

import (
  "net"
	"fmt"
)

func ReceiveDatagrams () (err error) {

	var conn *net.UDPConn

  if conn, err = createUDPListener("192.168.100.245:55555"); err != nil {
  	return err
  }

  fmt.Printf("Listener for UDP connections on %s\n", conn.LocalAddr().String())

  for {  	
  	rcv(conn)
  }

  panic("should never have got myself into this.")
}

func rcv (conn *net.UDPConn) {
		
		//func (c *UDPConn) ReadFromUDP(b []byte) (n int, addr *UDPAddr, err error)
  	
  	buffer := make([]byte, 8092)

  	if c, addr, err := conn.ReadFromUDP(buffer); err != nil {
  		
  		fmt.Println("blergh: " + err.Error())
  		return
  	
  	} else {
  		
  		fmt.Printf("datagram received from: %s\n", addr.String())
			fmt.Printf("bytes received: %d\n", c)
			fmt.Printf("buffer length: %d", len(buffer)) 
  	}		 	
}

func handleConnection (conn net.Conn) {
	fmt.Printf("datagram received from: %s\n", conn.RemoteAddr().String())
}

func createUDPListener (hostAndPort string) (conn *net.UDPConn, err error) {

	var udpaddr *net.UDPAddr
	if udpaddr, err = net.ResolveUDPAddr("udp", hostAndPort); err != nil {
		return
	}

  conn, err = net.ListenUDP("udp", udpaddr)
 
	return
}

