package main

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
	"fmt"
	"time"

	"github.com/bhojpur/api/pkg/websocket"
	"github.com/bhojpur/api/pkg/websocket/wsjson"
)

func main() {
	fmt.Printf("Bhojpur API - Echo Client trying to connect")
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err1 := websocket.Dial(ctx, "ws://localhost:8080", nil)
	if err1 != nil {
		fmt.Errorf("Bhojpur API echo server not available", err1)
		panic(err1)
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")
	fmt.Printf("connected successfully to Bhojpur API - Echo Server")

	err2 := wsjson.Write(ctx, c, "hi")
	if err2 != nil {
		fmt.Errorf("Bhojpur API unable to write JSON to echo server", err2)
		panic(err2)
	}
	fmt.Printf("written JSON to Bhojpur API - Echo Server")

	c.Close(websocket.StatusNormalClosure, "")

}
