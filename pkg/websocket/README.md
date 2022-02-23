# Bhojpur API - WebSocket Library

The `websocket` library provides __high-__ and __low-__ level bindings for the web browser's
WebSocket API.

## Installation

```bash
go get github.com/bhojpur/api
```

## Highlights

- Minimal and idiomatic API
- First class [context.Context](https://blog.golang.org/context) support
- Fully passes the WebSocket [autobahn-testsuite](https://github.com/crossbario/autobahn-testsuite)
- Single dependency
- JSON and protobuf helpers in the [wsjson](https://github.com/bhojpur/api/pkg/websocket/wsjson) and
 [wspb](https://github.com/bhojpur/api/pkg/websocket/wspb) sub-packages
- Zero alloc reads and writes
- Concurrent writes
- Close handshake
- net.Conn wrapper
- Ping pong API
- [RFC 7692](https://tools.ietf.org/html/rfc7692) permessage-deflate compression
- Compile to Wasm

## Simple Usage

### Server-side

```go
http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		// ...
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	var v interface{}
	err = wsjson.Read(ctx, c, &v)
	if err != nil {
		// ...
	}

	log.Printf("received: %v", v)

	c.Close(websocket.StatusNormalClosure, "")
})
```

### Client-side

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
defer cancel()

c, _, err := websocket.Dial(ctx, "ws://localhost:8080", nil)
if err != nil {
	// ...
}
defer c.Close(websocket.StatusInternalError, "the sky is falling")

err = wsjson.Write(ctx, c, "hi")
if err != nil {
	// ...
}

c.Close(websocket.StatusNormalClosure, "")
```

## Some Examples

For a production quality example that demonstrates the complete API, check the
[echo example](../../internal/echo).

For a full stack example, see the [chat example](../../internal/chat).
