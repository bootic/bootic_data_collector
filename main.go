package main

import (
	"datagram.io/daemon"
  "code.google.com/p/go.net/websocket"
	"datagram.io/ws"
  "net/http"
  "fmt"
  "text/template"
  "log"
)

const hostAndPort = "localhost:5555"

var homeTempl = template.Must(template.ParseFiles("views/home.html"))

func homeHandler(c http.ResponseWriter, req *http.Request) {
	homeTempl.Execute(c, req.Host)
}

func main () {
  
  go ws.Hub.Run()
  
	http.Handle("/ws", websocket.Handler(ws.WsHandler))
	fmt.Println("serving ws at " + hostAndPort + "/ws")
	
	http.HandleFunc("/", homeHandler)
	fmt.Println("serving HTTP at " + hostAndPort + "/")
	
	input := daemon.ReceiveDatagrams(hostAndPort)
	
	// Push incoming UDP messages to multiple listeners
	ws.Hub.Receive(input)
	
	fmt.Println("listening to UDP " + hostAndPort)
	
	if err := http.ListenAndServe(hostAndPort, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
	
}

func checkError (err error, info string) {
	if (err != nil) {
		panic("ERROR: " + info + ", " + err.Error())
	}
}
