package main

import (
	"datagram.io/daemon"
)

func main () {
	checkError(daemon.ReceiveDatagrams(), "Could not start UDP receive")
}

func checkError (err error, info string) {
	if (err != nil) {
		panic("ERROR: " + info + ", " + err.Error())
	}
}
