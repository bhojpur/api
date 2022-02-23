package wstest

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
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/bhojpur/api/pkg/websocket"
	"github.com/bhojpur/api/pkg/websocket/internal/test/assert"
	"github.com/bhojpur/api/pkg/websocket/internal/test/xrand"
	"github.com/bhojpur/api/pkg/websocket/internal/xsync"
)

// EchoLoop echos every msg received from c until an error
// occurs or the context expires.
// The read limit is set to 1 << 30.
func EchoLoop(ctx context.Context, c *websocket.Conn) error {
	defer c.Close(websocket.StatusInternalError, "")

	c.SetReadLimit(1 << 30)

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	b := make([]byte, 32<<10)
	for {
		typ, r, err := c.Reader(ctx)
		if err != nil {
			return err
		}

		w, err := c.Writer(ctx, typ)
		if err != nil {
			return err
		}

		_, err = io.CopyBuffer(w, r, b)
		if err != nil {
			return err
		}

		err = w.Close()
		if err != nil {
			return err
		}
	}
}

// Echo writes a message and ensures the same is sent back on c.
func Echo(ctx context.Context, c *websocket.Conn, max int) error {
	expType := websocket.MessageBinary
	if xrand.Bool() {
		expType = websocket.MessageText
	}

	msg := randMessage(expType, xrand.Int(max))

	writeErr := xsync.Go(func() error {
		return c.Write(ctx, expType, msg)
	})

	actType, act, err := c.Read(ctx)
	if err != nil {
		return err
	}

	err = <-writeErr
	if err != nil {
		return err
	}

	if expType != actType {
		return fmt.Errorf("unexpected message type (%v): %v", expType, actType)
	}

	if !bytes.Equal(msg, act) {
		return fmt.Errorf("unexpected msg read: %v", assert.Diff(msg, act))
	}

	return nil
}

func randMessage(typ websocket.MessageType, n int) []byte {
	if typ == websocket.MessageBinary {
		return xrand.Bytes(n)
	}
	return []byte(xrand.String(n))
}
