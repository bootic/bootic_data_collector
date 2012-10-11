package daemon

import (
  "net"
  "fmt"
)

type Input struct {
  Raw chan string
}

func (this *Input) Process(msg string) {
  this.Raw <- msg
}

func newInput() *Input {
  return &Input{
    Raw: make(chan string),
  }
}

func ReceiveDatagrams (hostAndPort string) Input {

	var conn *net.UDPConn
  
  conn, err := createUDPListener(hostAndPort)
  
  if err != nil {
  	panic("Could not create UDP listener")
  }

  fmt.Printf("Listener for UDP connections on %s\n", conn.LocalAddr().String())
  
  input := newInput()
  
  go rcv(conn, input)
  
  return *input
  
}

func rcv (conn *net.UDPConn, input *Input) {
	for {
	  buffer := make([]byte, 256)

  	if c, addr, err := conn.ReadFromUDP(buffer); err != nil {

  		fmt.Println("blergh: " + err.Error())
  		return

  	} else {
      
      msg := string(buffer[:c])
  		fmt.Printf("%d byte datagram received from %s\n\n", c, addr.String())
  		fmt.Printf("\t\"%s\"\n\n", msg)
  		
  		input.Process(msg)
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

