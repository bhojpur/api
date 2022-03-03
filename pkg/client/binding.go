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

	"github.com/pkg/errors"

	pb "github.com/bhojpur/application/pkg/api/v1/runtime"
)

// InvokeBindingRequest represents binding invocation request.
type InvokeBindingRequest struct {
	// Name is name of binding to invoke.
	Name string
	// Operation is the name of the operation type for the binding to invoke
	Operation string
	// Data is the input bindings sent
	Data []byte
	// Metadata is the input binding metadata
	Metadata map[string]string
}

// BindingEvent represents the binding event handler input.
type BindingEvent struct {
	// Data is the input bindings sent
	Data []byte
	// Metadata is the input binding metadata
	Metadata map[string]string
}

// InvokeBinding invokes specific operation on the configured Bhojpur Application runtime binding.
// This method covers input, output, and bi-directional bindings.
func (c *GRPCClient) InvokeBinding(ctx context.Context, in *InvokeBindingRequest) (*BindingEvent, error) {
	if in == nil {
		return nil, errors.New("binding invocation required")
	}
	if in.Name == "" {
		return nil, errors.New("binding invocation name required")
	}
	if in.Operation == "" {
		return nil, errors.New("binding invocation operation required")
	}

	req := &pb.InvokeBindingRequest{
		Name:      in.Name,
		Operation: in.Operation,
		Data:      in.Data,
		Metadata:  in.Metadata,
	}

	resp, err := c.protoClient.InvokeBinding(c.withAuthToken(ctx), req)
	if err != nil {
		return nil, errors.Wrapf(err, "error invoking binding %s/%s", in.Name, in.Operation)
	}

	if resp != nil {
		return &BindingEvent{
			Data:     resp.Data,
			Metadata: resp.Metadata,
		}, nil
	}

	return nil, nil
}

// InvokeOutputBinding invokes configured Bhojpur Application runtime binding with data (allows nil).InvokeOutputBinding
// This method differs from InvokeBinding in that it doesn't expect any content being returned from the invoked method.
func (c *GRPCClient) InvokeOutputBinding(ctx context.Context, in *InvokeBindingRequest) error {
	if _, err := c.InvokeBinding(ctx, in); err != nil {
		return errors.Wrap(err, "error invoking output binding")
	}
	return nil
}
