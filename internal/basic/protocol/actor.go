package protocol

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

import "context"

type ClientStub struct {
	GetUser         func(context.Context, *User) (*User, error)
	Invoke          func(context.Context, string) (string, error)
	Get             func(context.Context) (string, error)
	Post            func(context.Context, string) error
	StartTimer      func(context.Context, *TimerRequest) error
	StopTimer       func(context.Context, *TimerRequest) error
	StartReminder   func(context.Context, *ReminderRequest) error
	StopReminder    func(context.Context, *ReminderRequest) error
	IncrementAndGet func(ctx context.Context, stateKey string) (*User, error)
}

func (a *ClientStub) Type() string {
	return "testActorType"
}

func (a *ClientStub) ID() string {
	return "ActorImplID123456"
}

type User struct {
	Name string `json:"name"`
	Age  uint32 `json:"age"`
}

type TimerRequest struct {
	TimerName string `json:"timer_name"`
	CallBack  string `json:"call_back"`
	Duration  string `json:"duration"`
	Period    string `json:"period"`
	Data      string `json:"data"`
}

type ReminderRequest struct {
	ReminderName string `json:"reminder_name"`
	Duration     string `json:"duration"`
	Period       string `json:"period"`
	Data         string `json:"data"`
}
