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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	actorErr "github.com/bhojpur/api/pkg/actor/error"
	actorMock "github.com/bhojpur/api/pkg/actor/mock"
)

const mockActorID = "mockActorID"

func TestNewDefaultContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockServer := actorMock.NewMockServer(ctrl)
	mockCodec := actorMock.NewMockCodec(ctrl)

	mockServer.EXPECT().SetID(mockActorID)
	mockServer.EXPECT().SetStateManager(gomock.Any())
	mockServer.EXPECT().SaveState()
	mockServer.EXPECT().Type()

	newContainer, aerr := NewDefaultActorContainer(mockActorID, mockServer, mockCodec)
	assert.Equal(t, actorErr.Success, aerr)
	container, ok := newContainer.(*DefaultActorContainer)

	assert.True(t, ok)
	assert.NotNil(t, container)
	assert.NotNil(t, container.actor)
	assert.NotNil(t, container.serializer)
	assert.NotNil(t, container.methodType)
}

func TestContainerInvoke(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockServer := actorMock.NewMockServer(ctrl)
	mockCodec := actorMock.NewMockCodec(ctrl)
	param := `"param"`

	mockServer.EXPECT().SetID(mockActorID)
	mockServer.EXPECT().SetStateManager(gomock.Any())
	mockServer.EXPECT().SaveState()
	mockServer.EXPECT().Type()

	newContainer, aerr := NewDefaultActorContainer("mockActorID", mockServer, mockCodec)
	assert.Equal(t, actorErr.Success, aerr)
	container := newContainer.(*DefaultActorContainer)

	mockServer.EXPECT().Invoke(gomock.Any(), "param").Return(param, nil)
	mockCodec.EXPECT().Unmarshal([]byte(param), gomock.Any()).SetArg(1, "param").Return(nil)

	rsp, err := container.Invoke("Invoke", []byte(param))

	assert.Equal(t, 2, len(rsp))
	assert.Equal(t, actorErr.Success, err)
	assert.Equal(t, param, rsp[0].Interface().(string))
}
