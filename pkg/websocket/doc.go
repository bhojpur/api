//go:build !js
// +build !js

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

package websocket // import "github.com/bhojpur/api/pkg/websocket"

// It implements the WebSocket protocol (https://tools.ietf.org/html/rfc6455)
//
// Use Dial to dial a WebSocket server.
//
// Use Accept to accept a WebSocket client.
//
// Conn represents the resulting WebSocket connection.
//
// The examples are the best way to understand how to correctly use the library.
//
// The `wsjson` and `wspb` subpackages contain helpers for JSON and protobuf messages.
//
// Wasm
//
// The client-side supports compiling to Wasm. It wraps the WebSocket browser API.
//
// See https://developer.mozilla.org/en-US/docs/Web/API/WebSocket
//
// Some important caveats to be aware of:
//  - Accept always errors out
//  - conn.Ping is no-op
//  - HTTPClient, HTTPHeader and CompressionMode in DialOptions are no-op
//  - *http.Response from Dial is &http.Response{} with a 101 status code on success
