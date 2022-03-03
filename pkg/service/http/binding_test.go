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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/api/pkg/service/common"
)

func TestBindingHandlerWithoutHandler(t *testing.T) {
	s := newServer("", nil)
	err := s.AddBindingInvocationHandler("/", nil)
	assert.Errorf(t, err, "expected error adding nil binding event handler")
}

func TestBindingHandlerWithoutData(t *testing.T) {
	s := newServer("", nil)
	err := s.AddBindingInvocationHandler("/", func(ctx context.Context, in *common.BindingEvent) (out []byte, err error) {
		if in == nil {
			return nil, errors.New("nil input")
		}
		if in.Data != nil {
			return nil, errors.New("invalid input data")
		}
		return nil, nil
	})
	assert.NoErrorf(t, err, "error adding binding event handler")

	req, err := http.NewRequest(http.MethodPost, "/", nil)
	assert.NoErrorf(t, err, "error creating request")
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	s.mux.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "{}", resp.Body.String())
}

func TestBindingHandlerWithData(t *testing.T) {
	data := `{"name": "test"}`
	s := newServer("", nil)
	err := s.AddBindingInvocationHandler("/", func(ctx context.Context, in *common.BindingEvent) (out []byte, err error) {
		if in == nil {
			return nil, errors.New("nil input")
		}
		return []byte("test"), nil
	})
	assert.NoErrorf(t, err, "error adding binding event handler")

	req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(data))
	assert.NoErrorf(t, err, "error creating request")
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	s.mux.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "test", resp.Body.String())
}

func bindingHandlerFn(ctx context.Context, in *common.BindingEvent) (out []byte, err error) {
	if in == nil {
		return nil, errors.New("nil input")
	}
	return []byte("test"), nil
}

func bindingHandlerFnWithError(ctx context.Context, in *common.BindingEvent) (out []byte, err error) {
	return nil, errors.New("intentional error")
}

func TestBindingHandlerErrors(t *testing.T) {
	data := `{"name": "test"}`
	s := newServer("", nil)
	err := s.AddBindingInvocationHandler("", bindingHandlerFn)
	assert.Errorf(t, err, "expected error adding binding event handler sans route")

	err = s.AddBindingInvocationHandler("errors", bindingHandlerFnWithError)
	assert.NoErrorf(t, err, "error adding binding event handler sans slash")

	req, err := http.NewRequest(http.MethodPost, "/errors", strings.NewReader(data))
	assert.NoErrorf(t, err, "error creating request")
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	s.mux.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}
