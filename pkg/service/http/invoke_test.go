package http

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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/api/pkg/service/common"
)

func TestInvocationHandlerWithoutHandler(t *testing.T) {
	s := newServer("", nil)
	err := s.AddServiceInvocationHandler("/", nil)
	assert.Errorf(t, err, "expected error adding event handler")
}

func TestInvocationHandlerWithToken(t *testing.T) {
	data := `{"name": "test", "data": hellow}`
	_ = os.Setenv(common.AppAPITokenEnvVar, "app-app-token")
	s := newServer("", nil)
	err := s.AddServiceInvocationHandler("/", func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
		if in == nil || in.Data == nil || in.ContentType == "" {
			err = errors.New("nil input")
			return
		}
		out = &common.Content{
			Data:        in.Data,
			ContentType: in.ContentType,
			DataTypeURL: in.DataTypeURL,
		}
		return
	})
	assert.NoErrorf(t, err, "error adding event handler")

	// forbbiden.
	req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(data))
	assert.NoErrorf(t, err, "error creating request")
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	s.mux.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusNonAuthoritativeInfo, resp.Code)

	// pass.
	req.Header.Set(common.APITokenKey, os.Getenv(common.AppAPITokenEnvVar))
	resp = httptest.NewRecorder()
	s.mux.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	_ = os.Unsetenv(common.AppAPITokenEnvVar)
}

func TestInvocationHandlerWithData(t *testing.T) {
	data := `{"name": "test", "data": hellow}`
	s := newServer("", nil)
	err := s.AddServiceInvocationHandler("/", func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
		if in == nil || in.Data == nil || in.ContentType == "" {
			err = errors.New("nil input")
			return
		}
		out = &common.Content{
			Data:        in.Data,
			ContentType: in.ContentType,
			DataTypeURL: in.DataTypeURL,
		}
		return
	})
	assert.NoErrorf(t, err, "error adding event handler")

	req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(data))
	assert.NoErrorf(t, err, "error creating request")
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	s.mux.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	b, err := ioutil.ReadAll(resp.Body)
	assert.NoErrorf(t, err, "error reading response body")
	assert.Equal(t, data, string(b))
}

func TestInvocationHandlerWithoutInputData(t *testing.T) {
	s := newServer("", nil)
	err := s.AddServiceInvocationHandler("/", func(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
		if in == nil || in.Data != nil {
			err = errors.New("nil input")
			return
		}
		return &common.Content{}, nil
	})
	assert.NoErrorf(t, err, "error adding event handler")

	req, err := http.NewRequest(http.MethodPost, "/", nil)
	assert.NoErrorf(t, err, "error creating request")
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	s.mux.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	b, err := ioutil.ReadAll(resp.Body)
	assert.NoErrorf(t, err, "error reading response body")
	assert.NotNil(t, b)
	assert.Equal(t, "", string(b))
}

func emptyInvocationFn(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	return nil, nil
}

func TestInvocationHandlerWithInvalidRoute(t *testing.T) {
	s := newServer("", nil)

	err := s.AddServiceInvocationHandler("no-slash", emptyInvocationFn)
	assert.NoErrorf(t, err, "error adding no slash route event handler")

	err = s.AddServiceInvocationHandler("", emptyInvocationFn)
	assert.Errorf(t, err, "expected error from adding no route event handler")

	err = s.AddServiceInvocationHandler("/a", emptyInvocationFn)
	assert.NoErrorf(t, err, "error adding event handler")

	makeEventRequest(t, s, "/b", "", http.StatusNotFound)
}

func errorInvocationFn(ctx context.Context, in *common.InvocationEvent) (out *common.Content, err error) {
	return nil, errors.New("intentional test error")
}

func TestInvocationHandlerWithError(t *testing.T) {
	s := newServer("", nil)

	err := s.AddServiceInvocationHandler("/error", errorInvocationFn)
	assert.NoErrorf(t, err, "error adding error event handler")

	makeEventRequest(t, s, "/error", "", http.StatusInternalServerError)
}
