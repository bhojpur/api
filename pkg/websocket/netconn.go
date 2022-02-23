package websocket

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
	"io"
	"math"
	"net"
	"sync"
	"time"
)

// NetConn converts a *websocket.Conn into a net.Conn.
//
// It's for tunneling arbitrary protocols over WebSockets.
// Few users of the library will need this but it's tricky to implement
// correctly and so provided in the library.
//
// Every Write to the net.Conn will correspond to a message write of
// the given type on *websocket.Conn.
//
// The passed ctx bounds the lifetime of the net.Conn. If cancelled,
// all reads and writes on the net.Conn will be cancelled.
//
// If a message is read that is not of the correct type, the connection
// will be closed with StatusUnsupportedData and an error will be returned.
//
// Close will close the *websocket.Conn with StatusNormalClosure.
//
// When a deadline is hit, the connection will be closed. This is
// different from most net.Conn implementations where only the
// reading/writing goroutines are interrupted but the connection is kept alive.
//
// The Addr methods will return a mock net.Addr that returns "websocket" for Network
// and "websocket/unknown-addr" for String.
//
// A received StatusNormalClosure or StatusGoingAway close frame will be translated to
// io.EOF when reading.
func NetConn(ctx context.Context, c *Conn, msgType MessageType) net.Conn {
	nc := &netConn{
		c:       c,
		msgType: msgType,
	}

	var cancel context.CancelFunc
	nc.writeContext, cancel = context.WithCancel(ctx)
	nc.writeTimer = time.AfterFunc(math.MaxInt64, cancel)
	if !nc.writeTimer.Stop() {
		<-nc.writeTimer.C
	}

	nc.readContext, cancel = context.WithCancel(ctx)
	nc.readTimer = time.AfterFunc(math.MaxInt64, cancel)
	if !nc.readTimer.Stop() {
		<-nc.readTimer.C
	}

	return nc
}

type netConn struct {
	c       *Conn
	msgType MessageType

	writeTimer   *time.Timer
	writeContext context.Context

	readTimer   *time.Timer
	readContext context.Context

	readMu sync.Mutex
	eofed  bool
	reader io.Reader
}

var _ net.Conn = &netConn{}

func (c *netConn) Close() error {
	return c.c.Close(StatusNormalClosure, "")
}

func (c *netConn) Write(p []byte) (int, error) {
	err := c.c.Write(c.writeContext, c.msgType, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (c *netConn) Read(p []byte) (int, error) {
	c.readMu.Lock()
	defer c.readMu.Unlock()

	if c.eofed {
		return 0, io.EOF
	}

	if c.reader == nil {
		typ, r, err := c.c.Reader(c.readContext)
		if err != nil {
			switch CloseStatus(err) {
			case StatusNormalClosure, StatusGoingAway:
				c.eofed = true
				return 0, io.EOF
			}
			return 0, err
		}
		if typ != c.msgType {
			err := fmt.Errorf("unexpected frame type read (expected %v): %v", c.msgType, typ)
			c.c.Close(StatusUnsupportedData, err.Error())
			return 0, err
		}
		c.reader = r
	}

	n, err := c.reader.Read(p)
	if err == io.EOF {
		c.reader = nil
		err = nil
	}
	return n, err
}

type websocketAddr struct {
}

func (a websocketAddr) Network() string {
	return "websocket"
}

func (a websocketAddr) String() string {
	return "websocket/unknown-addr"
}

func (c *netConn) RemoteAddr() net.Addr {
	return websocketAddr{}
}

func (c *netConn) LocalAddr() net.Addr {
	return websocketAddr{}
}

func (c *netConn) SetDeadline(t time.Time) error {
	c.SetWriteDeadline(t)
	c.SetReadDeadline(t)
	return nil
}

func (c *netConn) SetWriteDeadline(t time.Time) error {
	if t.IsZero() {
		c.writeTimer.Stop()
	} else {
		c.writeTimer.Reset(t.Sub(time.Now()))
	}
	return nil
}

func (c *netConn) SetReadDeadline(t time.Time) error {
	if t.IsZero() {
		c.readTimer.Stop()
	} else {
		c.readTimer.Reset(t.Sub(time.Now()))
	}
	return nil
}
