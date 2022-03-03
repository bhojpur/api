package error

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

type ActorErr uint8

// TODO: the classification, handle and print log of error should be optimized.
const (
	Success                       = ActorErr(0)
	ErrActorTypeNotFound          = ActorErr(1)
	ErrRemindersParamsInvalid     = ActorErr(2)
	ErrActorMethodNoFound         = ActorErr(3)
	ErrActorInvokeFailed          = ActorErr(4)
	ErrReminderFuncUndefined      = ActorErr(5)
	ErrActorMethodSerializeFailed = ActorErr(6)
	ErrActorSerializeNoFound      = ActorErr(7)
	ErrActorIDNotFound            = ActorErr(8)
	ErrActorFactoryNotSet         = ActorErr(9)
	ErrTimerParamsInvalid         = ActorErr(10)
	ErrSaveStateFailed            = ActorErr(11)
	ErrActorServerInvalid         = ActorErr(12)
)
