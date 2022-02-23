package websocket_test

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/bhojpur/api/pkg/websocket"
	"github.com/bhojpur/api/pkg/websocket/wsjson"
)

func ExampleAccept() {
	// This handler accepts a WebSocket connection, reads a single JSON
	// message from the client and then closes the connection.

	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer c.Close(websocket.StatusInternalError, "the sky is falling")

		ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
		defer cancel()

		var v interface{}
		err = wsjson.Read(ctx, c, &v)
		if err != nil {
			log.Println(err)
			return
		}

		c.Close(websocket.StatusNormalClosure, "")
	})

	err := http.ListenAndServe("localhost:8080", fn)
	log.Fatal(err)
}

func ExampleDial() {
	// Dials a server, writes a single JSON message and then
	// closes the connection.

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	err = wsjson.Write(ctx, c, "hi")
	if err != nil {
		log.Fatal(err)
	}

	c.Close(websocket.StatusNormalClosure, "")
}

func ExampleCloseStatus() {
	// Dials a server and then expects to be disconnected with status code
	// websocket.StatusNormalClosure.

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	_, _, err = c.Reader(ctx)
	if websocket.CloseStatus(err) != websocket.StatusNormalClosure {
		log.Fatalf("expected to be disconnected with StatusNormalClosure but got: %v", err)
	}
}

func Example_writeOnly() {
	// This handler demonstrates how to correctly handle a write only WebSocket connection.
	// i.e you only expect to write messages and do not expect to read any messages.
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer c.Close(websocket.StatusInternalError, "the sky is falling")

		ctx, cancel := context.WithTimeout(r.Context(), time.Minute*10)
		defer cancel()

		ctx = c.CloseRead(ctx)

		t := time.NewTicker(time.Second * 30)
		defer t.Stop()

		for {
			select {
			case <-ctx.Done():
				c.Close(websocket.StatusNormalClosure, "")
				return
			case <-t.C:
				err = wsjson.Write(ctx, c, "hi")
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	})

	err := http.ListenAndServe("localhost:8080", fn)
	log.Fatal(err)
}

func Example_crossOrigin() {
	// This handler demonstrates how to safely accept cross origin WebSockets
	// from the origin bhojpur.net
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			OriginPatterns: []string{"bhojpur.net"},
		})
		if err != nil {
			log.Println(err)
			return
		}
		c.Close(websocket.StatusNormalClosure, "cross origin WebSocket accepted")
	})

	err := http.ListenAndServe("localhost:8080", fn)
	log.Fatal(err)
}

// This example demonstrates how to create a WebSocket server
// that gracefully exits when sent a signal.
//
// It starts a WebSocket server that keeps every connection open
// for 10 seconds.
// If you CTRL+C while a connection is open, it will wait at most 30s
// for all connections to terminate before shutting down.
// func ExampleGrace() {
// 	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		c, err := websocket.Accept(w, r, nil)
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}
// 		defer c.Close(websocket.StatusInternalError, "the sky is falling")
//
// 		ctx := c.CloseRead(r.Context())
// 		select {
// 		case <-ctx.Done():
// 		case <-time.After(time.Second * 10):
// 		}
//
// 		c.Close(websocket.StatusNormalClosure, "")
// 	})
//
// 	var g websocket.Grace
// 	s := &http.Server{
// 		Handler:      g.Handler(fn),
// 		ReadTimeout:  time.Second * 15,
// 		WriteTimeout: time.Second * 15,
// 	}
//
// 	errc := make(chan error, 1)
// 	go func() {
// 		errc <- s.ListenAndServe()
// 	}()
//
// 	sigs := make(chan os.Signal, 1)
// 	signal.Notify(sigs, os.Interrupt)
// 	select {
// 	case err := <-errc:
// 		log.Printf("failed to listen and serve: %v", err)
// 	case sig := <-sigs:
// 		log.Printf("terminating: %v", sig)
// 	}
//
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
// 	defer cancel()
// 	s.Shutdown(ctx)
// 	g.Shutdown(ctx)
// }

// This example demonstrates full stack chat with an automated test.
func Example_fullStackChat() {
	// https://github.com/bhojpur/api/tree/master/internal/chat
}

// This example demonstrates a echo server.
func Example_echo() {
	// https://github.com/bhojpur/api/tree/master/internal/echo
}
