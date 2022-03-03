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
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"github.com/bhojpur/api/pkg/actor"
	"github.com/bhojpur/api/pkg/actor/config"
	"github.com/bhojpur/api/pkg/actor/runtime"
	"github.com/bhojpur/api/pkg/service/common"
	"github.com/bhojpur/api/pkg/service/internal"
)

// NewService creates new Service.
func NewService(address string) common.Service {
	return newServer(address, nil)
}

// NewServiceWithMux creates new Service with existing http mux.
func NewServiceWithMux(address string, mux *mux.Router) common.Service {
	return newServer(address, mux)
}

func newServer(address string, router *mux.Router) *Server {
	if router == nil {
		router = mux.NewRouter()
	}
	return &Server{
		address: address,
		httpServer: &http.Server{
			Addr:    address,
			Handler: router,
		},
		mux:            router,
		topicRegistrar: make(internal.TopicRegistrar),
		authToken:      os.Getenv(common.AppAPITokenEnvVar),
	}
}

// Server is the HTTP server wrapping mux many Bhojpur Application runtime helpers.
type Server struct {
	address        string
	mux            *mux.Router
	httpServer     *http.Server
	topicRegistrar internal.TopicRegistrar
	authToken      string
}

func (s *Server) RegisterActorImplFactory(f actor.Factory, opts ...config.Option) {
	runtime.GetActorRuntimeInstance().RegisterActorFactory(f, opts...)
}

// Start starts the HTTP handler. Blocks while serving.
func (s *Server) Start() error {
	s.registerBaseHandler()
	return s.httpServer.ListenAndServe()
}

// Stop stops previously started HTTP service with a five second timeout.
func (s *Server) Stop() error {
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.httpServer.Shutdown(ctxShutDown)
}

func setOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
	w.Header().Set("Allow", "POST,OPTIONS")
}

func optionsHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			setOptions(w, r)
		} else {
			h.ServeHTTP(w, r)
		}
	}
}
