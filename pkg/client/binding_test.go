package client

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
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -timeout 30s ./client -count 1 -run ^TestInvokeBinding$

func TestInvokeBinding(t *testing.T) {
	ctx := context.Background()
	in := &InvokeBindingRequest{
		Name:      "test",
		Operation: "fn",
	}

	t.Run("output binding without data", func(t *testing.T) {
		err := testClient.InvokeOutputBinding(ctx, in)
		assert.Nil(t, err)
	})

	t.Run("output binding", func(t *testing.T) {
		in.Data = []byte("test")
		err := testClient.InvokeOutputBinding(ctx, in)
		assert.Nil(t, err)
	})

	t.Run("binding without data", func(t *testing.T) {
		in.Data = nil
		out, err := testClient.InvokeBinding(ctx, in)
		assert.Nil(t, err)
		assert.NotNil(t, out)
	})

	t.Run("binding with data and meta", func(t *testing.T) {
		in.Data = []byte("test")
		in.Metadata = map[string]string{"k1": "v1", "k2": "v2"}
		out, err := testClient.InvokeBinding(ctx, in)
		assert.Nil(t, err)
		assert.NotNil(t, out)
		assert.Equal(t, "test", string(out.Data))
	})
}
