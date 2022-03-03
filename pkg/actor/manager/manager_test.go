package manager

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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/api/pkg/actor/protocol"
	actorErr "github.com/bhojpur/api/pkg/actor/error"
	"github.com/bhojpur/api/pkg/actor/mock"
)

func TestNewDefaultActorManager(t *testing.T) {
	mng, err := NewDefaultActorManager("json")
	assert.NotNil(t, mng)
	assert.Equal(t, actorErr.Success, err)

	mng, err = NewDefaultActorManager("badSerializerType")
	assert.Nil(t, mng)
	assert.Equal(t, actorErr.ErrActorSerializeNoFound, err)
}

func TestRegisterActorImplFactory(t *testing.T) {
	mng, err := NewDefaultActorManager("json")
	assert.NotNil(t, mng)
	assert.Equal(t, actorErr.Success, err)
	assert.Nil(t, mng.(*DefaultActorManager).factory)
	mng.RegisterActorImplFactory(mock.ActorImplFactory)
	assert.NotNil(t, mng.(*DefaultActorManager).factory)
}

func TestInvokeMethod(t *testing.T) {
	mng, err := NewDefaultActorManager("json")
	assert.NotNil(t, mng)
	assert.Equal(t, actorErr.Success, err)
	assert.Nil(t, mng.(*DefaultActorManager).factory)

	data, err := mng.InvokeMethod("testActorID", "testMethodName", []byte(`"hello"`))
	assert.Nil(t, data)
	assert.Equal(t, actorErr.ErrActorFactoryNotSet, err)

	mng.RegisterActorImplFactory(mock.ActorImplFactory)
	assert.NotNil(t, mng.(*DefaultActorManager).factory)
	data, err = mng.InvokeMethod("testActorID", "mockMethod", []byte(`"hello"`))
	assert.Nil(t, data)
	assert.Equal(t, actorErr.ErrActorMethodNoFound, err)

	data, err = mng.InvokeMethod("testActorID", "Invoke", []byte(`"hello"`))
	assert.Equal(t, data, []byte(`"hello"`))
	assert.Equal(t, actorErr.Success, err)
}

func TestDeactivateActor(t *testing.T) {
	mng, err := NewDefaultActorManager("json")
	assert.NotNil(t, mng)
	assert.Equal(t, actorErr.Success, err)
	assert.Nil(t, mng.(*DefaultActorManager).factory)

	err = mng.DeactivateActor("testActorID")
	assert.Equal(t, actorErr.ErrActorIDNotFound, err)

	mng.RegisterActorImplFactory(mock.ActorImplFactory)
	assert.NotNil(t, mng.(*DefaultActorManager).factory)
	mng.InvokeMethod("testActorID", "Invoke", []byte(`"hello"`))

	err = mng.DeactivateActor("testActorID")
	assert.Equal(t, actorErr.Success, err)
}

func TestInvokeReminder(t *testing.T) {
	mng, err := NewDefaultActorManager("json")
	assert.NotNil(t, mng)
	assert.Equal(t, actorErr.Success, err)
	assert.Nil(t, mng.(*DefaultActorManager).factory)

	err = mng.InvokeReminder("testActorID", "testReminderName", []byte(`"hello"`))
	assert.Equal(t, actorErr.ErrActorFactoryNotSet, err)

	mng.RegisterActorImplFactory(mock.ActorImplFactory)
	assert.NotNil(t, mng.(*DefaultActorManager).factory)
	err = mng.InvokeReminder("testActorID", "testReminderName", []byte(`"hello"`))
	assert.Equal(t, actorErr.ErrRemindersParamsInvalid, err)

	reminderParam, _ := json.Marshal(&protocol.ActorReminderParams{
		Data:    []byte("hello"),
		DueTime: "5s",
		Period:  "6s",
	})
	err = mng.InvokeReminder("testActorID", "testReminderName", reminderParam)
	assert.Equal(t, actorErr.Success, err)
}

func TestInvokeTimer(t *testing.T) {
	mng, err := NewDefaultActorManager("json")
	assert.NotNil(t, mng)
	assert.Equal(t, actorErr.Success, err)
	assert.Nil(t, mng.(*DefaultActorManager).factory)

	err = mng.InvokeTimer("testActorID", "testTimerName", []byte(`"hello"`))
	assert.Equal(t, actorErr.ErrActorFactoryNotSet, err)

	mng.RegisterActorImplFactory(mock.ActorImplFactory)
	assert.NotNil(t, mng.(*DefaultActorManager).factory)
	err = mng.InvokeTimer("testActorID", "testTimerName", []byte(`"hello"`))
	assert.Equal(t, actorErr.ErrTimerParamsInvalid, err)

	timerParam, _ := json.Marshal(&protocol.ActorTimerParam{
		Data:     []byte("hello"),
		DueTime:  "5s",
		Period:   "6s",
		CallBack: "Invoke",
	})
	err = mng.InvokeTimer("testActorID", "testTimerName", timerParam)
	assert.Equal(t, actorErr.ErrActorMethodSerializeFailed, err)

	timerParam, _ = json.Marshal(&protocol.ActorTimerParam{
		Data:     []byte("hello"),
		DueTime:  "5s",
		Period:   "6s",
		CallBack: "NoSuchMethod",
	})
	err = mng.InvokeTimer("testActorID", "testTimerName", timerParam)
	assert.Equal(t, actorErr.ErrActorMethodNoFound, err)

	timerParam, _ = json.Marshal(&protocol.ActorTimerParam{
		Data:     []byte(`"hello"`),
		DueTime:  "5s",
		Period:   "6s",
		CallBack: "Invoke",
	})
	err = mng.InvokeTimer("testActorID", "testTimerName", timerParam)
	assert.Equal(t, actorErr.Success, err)
}
