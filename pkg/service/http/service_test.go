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
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStoppingUnstartedService(t *testing.T) {
	s := newServer("", nil)
	assert.NotNil(t, s)
	err := s.Stop()
	assert.NoError(t, err)
}

func TestStoppingStartedService(t *testing.T) {
	s := newServer(":3333", nil)
	assert.NotNil(t, s)

	go func() {
		if err := s.Start(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	// Wait for the server to start
	time.Sleep(200 * time.Millisecond)
	assert.NoError(t, s.Stop())
}

func TestStartingStoppedService(t *testing.T) {
	s := newServer(":3333", nil)
	assert.NotNil(t, s)
	stopErr := s.Stop()
	assert.NoError(t, stopErr)

	startErr := s.Start()
	assert.Error(t, startErr, "expected starting a stopped server to raise an error")
	assert.Equal(t, startErr.Error(), http.ErrServerClosed.Error())
}

func TestSettingOptions(t *testing.T) {
	req, err := http.NewRequest(http.MethodOptions, "/", nil)
	assert.NoErrorf(t, err, "error creating request")
	w := httptest.NewRecorder()
	setOptions(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.NotNil(t, resp)
	assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "POST,OPTIONS", resp.Header.Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "authorization, origin, content-type, accept", resp.Header.Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "POST,OPTIONS", resp.Header.Get("Allow"))
}

func testRequest(t *testing.T, s *Server, r *http.Request, expectedStatusCode int) {
	rr := httptest.NewRecorder()
	s.mux.ServeHTTP(rr, r)
	resp := rr.Result()
	defer resp.Body.Close()
	assert.NotNil(t, resp)
	assert.Equal(t, expectedStatusCode, resp.StatusCode)
}

func testRequestWithResponseBody(t *testing.T, s *Server, r *http.Request, expectedStatusCode int, expectedBody []byte) {
	rr := httptest.NewRecorder()
	s.mux.ServeHTTP(rr, r)
	rez := rr.Result()
	defer rez.Body.Close()
	rspBody, err := io.ReadAll(rez.Body)
	assert.Nil(t, err)
	assert.NotNil(t, rez)
	assert.Equal(t, expectedStatusCode, rez.StatusCode)
	assert.Equal(t, expectedBody, rspBody)
}
