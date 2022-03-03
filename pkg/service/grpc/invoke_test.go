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
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/anypb"

	cc "github.com/bhojpur/api/pkg/service/common"
	"github.com/bhojpur/application/pkg/api/v1/common"
)

func testInvokeHandler(ctx context.Context, in *cc.InvocationEvent) (out *cc.Content, err error) {
	if in == nil {
		return
	}
	out = &cc.Content{
		ContentType: in.ContentType,
		Data:        in.Data,
	}
	return
}

func testInvokeHandlerWithError(ctx context.Context, in *cc.InvocationEvent) (out *cc.Content, err error) {
	return nil, errors.New("test error")
}

func TestInvokeErrors(t *testing.T) {
	server := getTestServer()
	err := server.AddServiceInvocationHandler("", nil)
	assert.Error(t, err)
	err = server.AddServiceInvocationHandler("test", nil)
	assert.Error(t, err)
}

func TestInvokeWithToken(t *testing.T) {
	_ = os.Setenv(cc.AppAPITokenEnvVar, "app-app-token")
	server := getTestServer()
	startTestServer(server)
	methodName := "test"
	err := server.AddServiceInvocationHandler(methodName, testInvokeHandler)
	assert.Nil(t, err)
	t.Run("invoke with token, return success", func(t *testing.T) {
		grpcMetadata := metadata.New(map[string]string{
			cc.APITokenKey: os.Getenv(cc.AppAPITokenEnvVar),
		})
		ctx := metadata.NewIncomingContext(context.Background(), grpcMetadata)
		in := &common.InvokeRequest{Method: methodName}
		_, err := server.OnInvoke(ctx, in)
		assert.Nil(t, err)
	})
	t.Run("invoke with empty token, return failed", func(t *testing.T) {
		in := &common.InvokeRequest{Method: methodName}
		_, err := server.OnInvoke(context.Background(), in)
		assert.Error(t, err)
	})
	t.Run("invoke with mismatch token, return failed", func(t *testing.T) {
		grpcMetadata := metadata.New(map[string]string{
			cc.APITokenKey: "mismatch-token",
		})
		ctx := metadata.NewOutgoingContext(context.Background(), grpcMetadata)
		in := &common.InvokeRequest{Method: methodName}
		_, err := server.OnInvoke(ctx, in)
		assert.Error(t, err)
	})
	_ = os.Unsetenv(cc.AppAPITokenEnvVar)
}

// go test -timeout 30s ./service/grpc -count 1 -run ^TestInvoke$
func TestInvoke(t *testing.T) {
	methodName := "test"
	methodNameWithError := "error"
	ctx := context.Background()

	server := getTestServer()
	err := server.AddServiceInvocationHandler(methodName, testInvokeHandler)
	assert.Nil(t, err)

	err = server.AddServiceInvocationHandler(methodNameWithError, testInvokeHandlerWithError)
	assert.Nil(t, err)

	startTestServer(server)

	t.Run("invoke without request", func(t *testing.T) {
		_, err := server.OnInvoke(ctx, nil)
		assert.Error(t, err)
	})

	t.Run("invoke request with invalid method name", func(t *testing.T) {
		in := &common.InvokeRequest{Method: "invalid"}
		_, err := server.OnInvoke(ctx, in)
		assert.Error(t, err)
	})

	t.Run("invoke request without data", func(t *testing.T) {
		in := &common.InvokeRequest{Method: methodName}
		_, err := server.OnInvoke(ctx, in)
		assert.NoError(t, err)
	})

	t.Run("invoke request with data", func(t *testing.T) {
		data := "hello there"
		dataContentType := "text/plain"
		in := &common.InvokeRequest{Method: methodName}
		in.Data = &anypb.Any{Value: []byte(data)}
		in.ContentType = dataContentType
		out, err := server.OnInvoke(ctx, in)
		assert.NoError(t, err)
		assert.NotNil(t, out)
		assert.Equal(t, dataContentType, out.ContentType)
		assert.Equal(t, data, string(out.Data.Value))
	})

	t.Run("invoke request with error", func(t *testing.T) {
		data := "hello there"
		dataContentType := "text/plain"
		in := &common.InvokeRequest{Method: methodNameWithError}
		in.Data = &anypb.Any{Value: []byte(data)}
		in.ContentType = dataContentType
		_, err := server.OnInvoke(ctx, in)
		assert.Error(t, err)
	})

	stopTestServer(t, server)
}
