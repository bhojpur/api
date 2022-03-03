package mock

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
)

func ActorImplFactory() actor.Server {
	return &ActorImpl{}
}

type ActorImpl struct {
	actor.ServerImplBase
}

func (t *ActorImpl) Type() string {
	return "testActorType"
}

func (t *ActorImpl) Invoke(ctx context.Context, req string) (string, error) {
	return req, nil
}

func (t *ActorImpl) ReminderCall(reminderName string, state []byte, dueTime string, period string) {
}

func NotReminderCalleeActorFactory() actor.Server {
	return &NotReminderCalleeActor{}
}

type NotReminderCalleeActor struct {
	actor.ServerImplBase
}

func (t *NotReminderCalleeActor) Type() string {
	return "testActorNotReminderCalleeType"
}
