// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/bhojpur/api/pkg/actor/manager (interfaces: ActorContainer)

// Package mock is a generated GoMock package.
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
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	actor "github.com/bhojpur/api/pkg/actor"
	error "github.com/bhojpur/api/pkg/actor/error"
)

// MockActorContainer is a mock of ActorContainer interface.
type MockActorContainer struct {
	ctrl     *gomock.Controller
	recorder *MockActorContainerMockRecorder
}

// MockActorContainerMockRecorder is the mock recorder for MockActorContainer.
type MockActorContainerMockRecorder struct {
	mock *MockActorContainer
}

// NewMockActorContainer creates a new mock instance.
func NewMockActorContainer(ctrl *gomock.Controller) *MockActorContainer {
	mock := &MockActorContainer{ctrl: ctrl}
	mock.recorder = &MockActorContainerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockActorContainer) EXPECT() *MockActorContainerMockRecorder {
	return m.recorder
}

// GetActor mocks base method.
func (m *MockActorContainer) GetActor() actor.Server {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActor")
	ret0, _ := ret[0].(actor.Server)
	return ret0
}

// GetActor indicates an expected call of GetActor.
func (mr *MockActorContainerMockRecorder) GetActor() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActor", reflect.TypeOf((*MockActorContainer)(nil).GetActor))
}

// Invoke mocks base method.
func (m *MockActorContainer) Invoke(arg0 string, arg1 []byte) ([]reflect.Value, error.ActorErr) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Invoke", arg0, arg1)
	ret0, _ := ret[0].([]reflect.Value)
	ret1, _ := ret[1].(error.ActorErr)
	return ret0, ret1
}

// Invoke indicates an expected call of Invoke.
func (mr *MockActorContainerMockRecorder) Invoke(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Invoke", reflect.TypeOf((*MockActorContainer)(nil).Invoke), arg0, arg1)
}