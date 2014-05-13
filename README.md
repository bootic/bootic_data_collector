# Data collector

This is a realtime events/statistics recording and distribution backend written in Go.
The backend is designed to receive messages via UDP and publish them via Websockets and a ZMQ pub/sub socket:

Incoming events/stats are broadcast out to any listening websocket clients.

## Start

    $ go run main.go --zmqsocket=tcp://:6000 --udphost=host:5555 --wshost=host:5555 --ssehost=localhost:5556 --accesstoken=foobar


### Options

Option        | Description                                       | example
------------- | ------------------------------------------------- | -----------------
`zmqsocket`   | ZMQ socket to publish events to.                  | `tcp://6000`
`udphost`     | UDP host:port to listen for incoming events on    | `someserver:5555`
`wshost`      | host:port for Websockets endpoint                 | `somepublicserver.com:80`
`ssehost`     | host:port for Server Sent Events endpoint         | `somepublicserver:5555`
`accesstoken` | Access token to authenticate public endpoints     | "foobar123"
