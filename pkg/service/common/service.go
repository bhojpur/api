package common

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

	"github.com/bhojpur/api/pkg/actor"
	"github.com/bhojpur/api/pkg/actor/config"
)

const (
	// AppAPITokenEnvVar is the environment variable for Bhojpur Application API token.
	AppAPITokenEnvVar = "APP_API_TOKEN" /* #nosec */
	APITokenKey       = "app-api-token" /* #nosec */
)

// Service represents Bhojpur Application callback service.
type Service interface {
	// AddServiceInvocationHandler appends provided service invocation handler with its name to the service.
	AddServiceInvocationHandler(name string, fn ServiceInvocationHandler) error
	// AddTopicEventHandler appends provided event handler with its topic and optional metadata to the service.
	// Note, retries are only considered when there is an error. Lack of error is considered as a success
	AddTopicEventHandler(sub *Subscription, fn TopicEventHandler) error
	// AddBindingInvocationHandler appends provided binding invocation handler with its name to the service.
	AddBindingInvocationHandler(name string, fn BindingInvocationHandler) error
	// RegisterActorImplFactory Register a new actor to actor runtime of go sdk
	RegisterActorImplFactory(f actor.Factory, opts ...config.Option)
	// Start starts service.
	Start() error
	// Stop stops the previously started service.
	Stop() error
}

type (
	ServiceInvocationHandler func(ctx context.Context, in *InvocationEvent) (out *Content, err error)
	TopicEventHandler        func(ctx context.Context, e *TopicEvent) (retry bool, err error)
	BindingInvocationHandler func(ctx context.Context, in *BindingEvent) (out []byte, err error)
)
