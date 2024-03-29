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
	"testing"

	"github.com/stretchr/testify/assert"
)

const testActorType = "test"

func TestInvokeActor(t *testing.T) {
	ctx := context.Background()
	in := &InvokeActorRequest{
		ActorID:   "fn",
		Method:    "mockMethod",
		Data:      []byte(`{hello}`),
		ActorType: testActorType,
	}

	t.Run("invoke actor without data ", func(t *testing.T) {
		in.Data = nil
		out, err := testClient.InvokeActor(ctx, in)
		in.Data = []byte(`{hello}`)
		assert.NoError(t, err)
		assert.NotNil(t, out)
	})

	t.Run("invoke actor without method", func(t *testing.T) {
		in.Method = ""
		out, err := testClient.InvokeActor(ctx, in)
		in.Method = "mockMethod"
		assert.Error(t, err)
		assert.Nil(t, out)
	})

	t.Run("invoke actor without id ", func(t *testing.T) {
		in.ActorID = ""
		out, err := testClient.InvokeActor(ctx, in)
		in.ActorID = "fn"
		assert.Error(t, err)
		assert.Nil(t, out)
	})

	t.Run("invoke actor without type", func(t *testing.T) {
		in.ActorType = ""
		out, err := testClient.InvokeActor(ctx, in)
		in.ActorType = testActorType
		assert.Error(t, err)
		assert.Nil(t, out)
	})

	t.Run("invoke actor without empty input", func(t *testing.T) {
		in = nil
		out, err := testClient.InvokeActor(ctx, in)
		assert.Error(t, err)
		assert.Nil(t, out)
	})
}

func TestRegisterActorReminder(t *testing.T) {
	ctx := context.Background()
	in := &RegisterActorReminderRequest{
		ActorID:   "fn",
		Data:      []byte(`{hello}`),
		ActorType: testActorType,
		Name:      "mockName",
		Period:    "2s",
		DueTime:   "4s",
		TTL:       "20s",
	}

	t.Run("invoke register actor reminder without actorType", func(t *testing.T) {
		in.ActorType = ""
		err := testClient.RegisterActorReminder(ctx, in)
		in.ActorType = testActorType
		assert.Error(t, err)
	})

	t.Run("invoke register actor reminder without id ", func(t *testing.T) {
		in.ActorID = ""
		err := testClient.RegisterActorReminder(ctx, in)
		in.ActorID = "fn"
		assert.Error(t, err)
	})

	t.Run("invoke register actor reminder without Name ", func(t *testing.T) {
		in.Name = ""
		err := testClient.RegisterActorReminder(ctx, in)
		in.Name = "mockName"
		assert.Error(t, err)
	})

	t.Run("invoke register actor reminder without period ", func(t *testing.T) {
		in.Period = ""
		err := testClient.RegisterActorReminder(ctx, in)
		in.Period = "2s"
		assert.NoError(t, err)
	})

	t.Run("invoke register actor reminder without dutTime ", func(t *testing.T) {
		in.DueTime = ""
		err := testClient.RegisterActorReminder(ctx, in)
		in.DueTime = "2s"
		assert.NoError(t, err)
	})

	t.Run("invoke register actor reminder without TTL ", func(t *testing.T) {
		in.TTL = ""
		err := testClient.RegisterActorReminder(ctx, in)
		in.TTL = "20s"
		assert.NoError(t, err)
	})

	t.Run("invoke register actor reminder ", func(t *testing.T) {
		assert.NoError(t, testClient.RegisterActorReminder(ctx, in))
	})

	t.Run("invoke register actor reminder with empty param", func(t *testing.T) {
		assert.Error(t, testClient.RegisterActorReminder(ctx, nil))
	})
}

func TestRegisterActorTimer(t *testing.T) {
	ctx := context.Background()
	in := &RegisterActorTimerRequest{
		ActorID:   "fn",
		Data:      []byte(`{hello}`),
		ActorType: testActorType,
		Name:      "mockName",
		Period:    "2s",
		DueTime:   "4s",
		TTL:       "20s",
		CallBack:  "mockFunc",
	}

	t.Run("invoke register actor timer without actorType", func(t *testing.T) {
		in.ActorType = ""
		err := testClient.RegisterActorTimer(ctx, in)
		in.ActorType = testActorType
		assert.Error(t, err)
	})

	t.Run("invoke register actor timer without id ", func(t *testing.T) {
		in.ActorID = ""
		err := testClient.RegisterActorTimer(ctx, in)
		in.ActorID = "fn"
		assert.Error(t, err)
	})

	t.Run("invoke register actor timer without Name ", func(t *testing.T) {
		in.Name = ""
		err := testClient.RegisterActorTimer(ctx, in)
		in.Name = "mockName"
		assert.Error(t, err)
	})

	t.Run("invoke register actor timer without period ", func(t *testing.T) {
		in.Period = ""
		err := testClient.RegisterActorTimer(ctx, in)
		in.Period = "2s"
		assert.NoError(t, err)
	})

	t.Run("invoke register actor timer without dutTime ", func(t *testing.T) {
		in.DueTime = ""
		err := testClient.RegisterActorTimer(ctx, in)
		in.DueTime = "4s"
		assert.NoError(t, err)
	})

	t.Run("invoke register actor timer without TTL ", func(t *testing.T) {
		in.TTL = ""
		err := testClient.RegisterActorTimer(ctx, in)
		in.TTL = "20s"
		assert.NoError(t, err)
	})

	t.Run("invoke register actor timer without callBack ", func(t *testing.T) {
		in.CallBack = ""
		err := testClient.RegisterActorTimer(ctx, in)
		in.CallBack = "mockFunc"
		assert.Error(t, err)
	})

	t.Run("invoke register actor timer without data ", func(t *testing.T) {
		in.Data = nil
		err := testClient.RegisterActorTimer(ctx, in)
		in.Data = []byte(`{hello}`)
		assert.NoError(t, err)
	})

	t.Run("invoke register actor timer", func(t *testing.T) {
		assert.NoError(t, testClient.RegisterActorTimer(ctx, in))
	})

	t.Run("invoke register actor timer with empty param", func(t *testing.T) {
		assert.Error(t, testClient.RegisterActorTimer(ctx, nil))
	})
}

func TestUnregisterActorReminder(t *testing.T) {
	ctx := context.Background()
	in := &UnregisterActorReminderRequest{
		ActorID:   "fn",
		ActorType: testActorType,
		Name:      "mockName",
	}

	t.Run("invoke unregister actor reminder without actorType", func(t *testing.T) {
		in.ActorType = ""
		err := testClient.UnregisterActorReminder(ctx, in)
		in.ActorType = testActorType
		assert.Error(t, err)
	})

	t.Run("invoke unregister actor reminder without id ", func(t *testing.T) {
		in.ActorID = ""
		err := testClient.UnregisterActorReminder(ctx, in)
		in.ActorID = "fn"
		assert.Error(t, err)
	})

	t.Run("invoke unregister actor reminder without Name ", func(t *testing.T) {
		in.Name = ""
		err := testClient.UnregisterActorReminder(ctx, in)
		in.Name = "mockName"
		assert.Error(t, err)
	})

	t.Run("invoke unregister actor reminder without period ", func(t *testing.T) {
		in.ActorType = ""
		err := testClient.UnregisterActorReminder(ctx, in)
		in.ActorType = testActorType
		assert.Error(t, err)
	})

	t.Run("invoke unregister actor reminder ", func(t *testing.T) {
		assert.NoError(t, testClient.UnregisterActorReminder(ctx, in))
	})

	t.Run("invoke unregister actor reminder with empty param", func(t *testing.T) {
		assert.Error(t, testClient.UnregisterActorReminder(ctx, nil))
	})
}

func TestRenameActorReminder(t *testing.T) {
	ctx := context.Background()

	registerReminderReq := &RegisterActorReminderRequest{
		ActorID:   "fn",
		Data:      []byte(`{hello}`),
		ActorType: testActorType,
		Name:      "oldName",
		Period:    "2s",
		DueTime:   "4s",
		TTL:       "20s",
	}

	testClient.RegisterActorReminder(ctx, registerReminderReq)

	renameReminderReq := &RenameActorReminderRequest{
		ActorID:   "fn",
		ActorType: testActorType,
		OldName:   "oldName",
		NewName:   "newName",
	}

	t.Run("invoke rename actor reminder without actorType", func(t *testing.T) {
		renameReminderReq.ActorType = ""
		err := testClient.RenameActorReminder(ctx, renameReminderReq)
		renameReminderReq.ActorType = testActorType
		assert.Error(t, err)
	})

	t.Run("invoke rename actor reminder without id ", func(t *testing.T) {
		renameReminderReq.ActorID = ""
		err := testClient.RenameActorReminder(ctx, renameReminderReq)
		renameReminderReq.ActorID = "fn"
		assert.Error(t, err)
	})

	t.Run("invoke rename actor reminder without oldName ", func(t *testing.T) {
		renameReminderReq.OldName = ""
		err := testClient.RenameActorReminder(ctx, renameReminderReq)
		renameReminderReq.OldName = "oldName"
		assert.Error(t, err)
	})

	t.Run("invoke rename actor reminder without newName ", func(t *testing.T) {
		renameReminderReq.NewName = ""
		err := testClient.RenameActorReminder(ctx, renameReminderReq)
		renameReminderReq.NewName = "newName"
		assert.Error(t, err)
	})

	t.Run("invoke rename actor reminder ", func(t *testing.T) {
		assert.NoError(t, testClient.RenameActorReminder(ctx, renameReminderReq))
	})

	t.Run("invoke rename actor reminder with empty param", func(t *testing.T) {
		assert.Error(t, testClient.RenameActorReminder(ctx, nil))
	})
}

func TestUnregisterActorTimer(t *testing.T) {
	ctx := context.Background()
	in := &UnregisterActorTimerRequest{
		ActorID:   "fn",
		ActorType: testActorType,
		Name:      "mockName",
	}

	t.Run("invoke unregister actor timer without actorType", func(t *testing.T) {
		in.ActorType = ""
		err := testClient.UnregisterActorTimer(ctx, in)
		in.ActorType = testActorType
		assert.Error(t, err)
	})

	t.Run("invoke register actor timer without id ", func(t *testing.T) {
		in.ActorID = ""
		err := testClient.UnregisterActorTimer(ctx, in)
		in.ActorID = "fn"
		assert.Error(t, err)
	})

	t.Run("invoke register actor timer without Name ", func(t *testing.T) {
		in.Name = ""
		err := testClient.UnregisterActorTimer(ctx, in)
		in.Name = "mockName"
		assert.Error(t, err)
	})

	t.Run("invoke register actor timer without period ", func(t *testing.T) {
		in.ActorType = ""
		err := testClient.UnregisterActorTimer(ctx, in)
		in.ActorType = testActorType
		assert.Error(t, err)
	})

	t.Run("invoke register actor timer ", func(t *testing.T) {
		assert.NoError(t, testClient.UnregisterActorTimer(ctx, in))
	})

	t.Run("invoke register actor timer with empty param", func(t *testing.T) {
		assert.Error(t, testClient.UnregisterActorTimer(ctx, nil))
	})
}
