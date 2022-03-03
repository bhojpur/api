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
	"net"
	"os"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	pb "github.com/bhojpur/application/pkg/api/v1/runtime"

	"github.com/bhojpur/api/pkg/actor"
	"github.com/bhojpur/api/pkg/actor/config"
	"github.com/bhojpur/api/pkg/service/common"
	"github.com/bhojpur/api/pkg/service/internal"
)

// NewService creates new Service.
func NewService(address string) (s common.Service, err error) {
	if address == "" {
		return nil, errors.New("nil address")
	}
	lis, err := net.Listen("tcp", address)
	if err != nil {
		err = errors.Wrapf(err, "failed to TCP listen on: %s", address)
		return
	}
	s = newService(lis)
	return
}

// NewServiceWithListener creates new Service with specific listener.
func NewServiceWithListener(lis net.Listener) common.Service {
	return newService(lis)
}

func newService(lis net.Listener) *Server {
	return &Server{
		listener:        lis,
		invokeHandlers:  make(map[string]common.ServiceInvocationHandler),
		topicRegistrar:  make(internal.TopicRegistrar),
		bindingHandlers: make(map[string]common.BindingInvocationHandler),
		authToken:       os.Getenv(common.AppAPITokenEnvVar),
	}
}

// Server is the gRPC service implementation for Bhojpur Application runtime.
type Server struct {
	pb.UnimplementedAppCallbackServer
	listener        net.Listener
	invokeHandlers  map[string]common.ServiceInvocationHandler
	topicRegistrar  internal.TopicRegistrar
	bindingHandlers map[string]common.BindingInvocationHandler
	authToken       string
}

func (s *Server) RegisterActorImplFactory(f actor.Factory, opts ...config.Option) {
	panic("Actor is not supported by gRPC API")
}

// Start registers the server and starts it.
func (s *Server) Start() error {
	gs := grpc.NewServer()
	pb.RegisterAppCallbackServer(gs, s)
	return gs.Serve(s.listener)
}

// Stop stops the previously started service.
func (s *Server) Stop() error {
	return s.listener.Close()
}
