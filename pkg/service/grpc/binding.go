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
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"

	"github.com/bhojpur/api/pkg/service/common"
	pb "github.com/bhojpur/application/pkg/api/v1/runtime"
)

// AddBindingInvocationHandler appends provided binding invocation handler with its name to the service.
func (s *Server) AddBindingInvocationHandler(name string, fn common.BindingInvocationHandler) error {
	if name == "" {
		return fmt.Errorf("binding name required")
	}
	if fn == nil {
		return fmt.Errorf("binding handler required")
	}
	s.bindingHandlers[name] = fn
	return nil
}

// ListInputBindings is called by Bhojpur Application runtime to get the list of bindings the application will get invoked by. In this example, we are telling Bhojpur Application runtime
// To invoke our application with a binding named storage.
func (s *Server) ListInputBindings(ctx context.Context, in *empty.Empty) (*pb.ListInputBindingsResponse, error) {
	list := make([]string, 0)
	for k := range s.bindingHandlers {
		list = append(list, k)
	}

	return &pb.ListInputBindingsResponse{
		Bindings: list,
	}, nil
}

// OnBindingEvent gets invoked every time a new event is fired from a registered binding. The message carries the binding name, a payload and optional metadata.
func (s *Server) OnBindingEvent(ctx context.Context, in *pb.BindingEventRequest) (*pb.BindingEventResponse, error) {
	if in == nil {
		return nil, errors.New("nil binding event request")
	}
	if fn, ok := s.bindingHandlers[in.Name]; ok {
		e := &common.BindingEvent{
			Data:     in.Data,
			Metadata: in.Metadata,
		}
		data, err := fn(ctx, e)
		if err != nil {
			return nil, errors.Wrapf(err, "error executing %s binding", in.Name)
		}
		return &pb.BindingEventResponse{
			Data: data,
		}, nil
	}

	return nil, fmt.Errorf("binding not implemented: %s", in.Name)
}
