# datagram.io

Datagram is a realtime events/statistics recording and distribution backend written in Go.
The backend is designed to receive messages via UDP and publish them via Websockets and a ZMQ pub/sub socket:

Incoming events/stats are broadcast out to any listening websocket clients.

## Start

    $ go run main.go --zmqsocket=tcp://:6000 --udphost=tracking_host:5555 --wshost=ws_host:5555
