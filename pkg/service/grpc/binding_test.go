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
	"context"
	"errors"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/api/pkg/service/common"
	runtime "github.com/bhojpur/application/pkg/api/v1/runtime"
)

func testBindingHandler(ctx context.Context, in *common.BindingEvent) (out []byte, err error) {
	if in == nil {
		return nil, errors.New("nil event")
	}
	return in.Data, nil
}

func TestListInputBindings(t *testing.T) {
	server := getTestServer()
	err := server.AddBindingInvocationHandler("test1", testBindingHandler)
	assert.NoError(t, err)
	err = server.AddBindingInvocationHandler("test2", testBindingHandler)
	assert.NoError(t, err)
	resp, err := server.ListInputBindings(context.Background(), &empty.Empty{})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Lenf(t, resp.Bindings, 2, "expected 2 handlers")
}

func TestBindingForErrors(t *testing.T) {
	server := getTestServer()
	err := server.AddBindingInvocationHandler("", nil)
	assert.Errorf(t, err, "expected error on nil method name")

	err = server.AddBindingInvocationHandler("test", nil)
	assert.Errorf(t, err, "expected error on nil method handler")
}

// go test -timeout 30s ./service/grpc -count 1 -run ^TestBinding$
func TestBinding(t *testing.T) {
	ctx := context.Background()
	methodName := "test"

	server := getTestServer()
	err := server.AddBindingInvocationHandler(methodName, testBindingHandler)
	assert.Nil(t, err)
	startTestServer(server)

	t.Run("binding without event", func(t *testing.T) {
		_, err := server.OnBindingEvent(ctx, nil)
		assert.Error(t, err)
	})

	t.Run("binding event for wrong method", func(t *testing.T) {
		in := &runtime.BindingEventRequest{Name: "invalid"}
		_, err := server.OnBindingEvent(ctx, in)
		assert.Error(t, err)
	})

	t.Run("binding event without data", func(t *testing.T) {
		in := &runtime.BindingEventRequest{Name: methodName}
		out, err := server.OnBindingEvent(ctx, in)
		assert.NoError(t, err)
		assert.NotNil(t, out)
	})

	t.Run("binding event with data", func(t *testing.T) {
		data := "hello there"
		in := &runtime.BindingEventRequest{
			Name: methodName,
			Data: []byte(data),
		}
		out, err := server.OnBindingEvent(ctx, in)
		assert.NoError(t, err)
		assert.NotNil(t, out)
		assert.Equal(t, data, string(out.Data))
	})

	t.Run("binding event with metadata", func(t *testing.T) {
		in := &runtime.BindingEventRequest{
			Name:     methodName,
			Metadata: map[string]string{"k1": "v1", "k2": "v2"},
		}
		out, err := server.OnBindingEvent(ctx, in)
		assert.NoError(t, err)
		assert.NotNil(t, out)
	})

	stopTestServer(t, server)
}
