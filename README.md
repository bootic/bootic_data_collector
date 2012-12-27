# datagram.io

Datagram is a realtime events/statistics recording and distribution backend written in Go (for performance reasons).
The backend is designed to receive messages via:

* UDP for optimal performance
* TCP for reliability
* HTTP for web integration

Incoming events/stats are broadcast out to any listening websocket clients. They are also stored in a postgres instance which uses a low-level postgres adapter for optimal speed.

An api to query for historical events is also available over HTTP.

The system uses postgres with carefully designed indicies as its primary datastore. This can be backed up by a Redis instance if/when postgres isn't able to handle the load.

## Start

    $ go run main.go daemons
