package grpc

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
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/test/bufconn"
)

func TestServer(t *testing.T) {
	server := getTestServer()
	startTestServer(server)
	stopTestServer(t, server)
}

func TestServerWithListener(t *testing.T) {
	server := NewServiceWithListener(bufconn.Listen(1024 * 1024))
	assert.NotNil(t, server)
}

func TestService(t *testing.T) {
	_, err := NewService("")
	assert.Errorf(t, err, "expected error from lack of address")
}

func getTestServer() *Server {
	return newService(bufconn.Listen(1024 * 1024))
}

func startTestServer(server *Server) {
	go func() {
		if err := server.Start(); err != nil && err.Error() != "closed" {
			panic(err)
		}
	}()
}

func stopTestServer(t *testing.T, server *Server) {
	assert.NotNil(t, server)
	err := server.Stop()
	assert.Nilf(t, err, "error stopping server")
}
