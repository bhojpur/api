# Bhojpur API - Echo Example

This directory contains an Echo server example using
[Bhojpur API - WebSocket](https://github.com/bhojpur/api/pkg/websocket)..

```bash
$ cd internal/echo
$ go run . localhost:3000
listening on http://127.0.0.1:3000
```

You can use a WebSocket client like [https://github.com/hashrocket/ws](https://github.com/hashrocket/ws)
to connect. All messages written will be echoed back.

## Structure

The server is in `server.go` and is implemented as a `http.HandlerFunc` that accepts the WebSocket
and then reads all messages and writes them exactly as is back to the connection.

`server_test.go` contains a small unit test to verify it works correctly.

`main.go` brings it all together so that you can run it and play around with it.