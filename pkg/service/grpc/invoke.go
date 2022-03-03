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

	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"

	cc "github.com/bhojpur/api/pkg/service/common"
	cpb "github.com/bhojpur/application/pkg/api/v1/common"
)

// AddServiceInvocationHandler appends provided service invocation handler with its method to the service.
func (s *Server) AddServiceInvocationHandler(method string, fn cc.ServiceInvocationHandler) error {
	if method == "" {
		return fmt.Errorf("service name required")
	}
	if fn == nil {
		return fmt.Errorf("invocation handler required")
	}
	s.invokeHandlers[method] = fn
	return nil
}

// OnInvoke gets invoked when a remote service has called the application through Bhojpur Application runtime.
func (s *Server) OnInvoke(ctx context.Context, in *cpb.InvokeRequest) (*cpb.InvokeResponse, error) {
	if in == nil {
		return nil, errors.New("nil invoke request")
	}
	if s.authToken != "" {
		if md, ok := metadata.FromIncomingContext(ctx); !ok {
			return nil, errors.New("authentication failed")
		} else if vals := md.Get(cc.APITokenKey); len(vals) > 0 {
			if vals[0] != s.authToken {
				return nil, errors.New("authentication failed: application token mismatch")
			}
		} else {
			return nil, errors.New("authentication failed. application token key does not exist")
		}
	}
	if fn, ok := s.invokeHandlers[in.Method]; ok {
		e := &cc.InvocationEvent{}
		e.ContentType = in.ContentType

		if in.Data != nil {
			e.Data = in.Data.Value
			e.DataTypeURL = in.Data.TypeUrl
		}

		if in.HttpExtension != nil {
			e.Verb = in.HttpExtension.Verb.String()
			e.QueryString = in.HttpExtension.Querystring
		}

		ct, er := fn(ctx, e)
		if er != nil {
			return nil, er
		}

		if ct == nil {
			return &cpb.InvokeResponse{}, nil
		}

		return &cpb.InvokeResponse{
			ContentType: ct.ContentType,
			Data: &any.Any{
				Value:   ct.Data,
				TypeUrl: ct.DataTypeURL,
			},
		}, nil
	}
	return nil, fmt.Errorf("method not implemented: %s", in.Method)
}
