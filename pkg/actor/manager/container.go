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
	"context"
	"log"
	"reflect"

	"github.com/bhojpur/api/pkg/actor"
	"github.com/bhojpur/api/pkg/actor/codec"
	actorErr "github.com/bhojpur/api/pkg/actor/error"
	"github.com/bhojpur/api/pkg/actor/state"
	app "github.com/bhojpur/api/pkg/client"
)

type ActorContainer interface {
	Invoke(methodName string, param []byte) ([]reflect.Value, actorErr.ActorErr)
	GetActor() actor.Server
}

// DefaultActorContainer contains actor instance and methods type info generated from actor.
type DefaultActorContainer struct {
	methodType map[string]*MethodType
	actor      actor.Server
	serializer codec.Codec
}

// NewDefaultActorContainer creates a new ActorContainer with provider impl actor and serializer.
func NewDefaultActorContainer(actorID string, impl actor.Server, serializer codec.Codec) (ActorContainer, actorErr.ActorErr) {
	impl.SetID(actorID)
	appClient, _ := app.NewClient()
	// create state manager for this new actor
	impl.SetStateManager(state.NewActorStateManager(impl.Type(), actorID, state.NewAppStateAsyncProvider(appClient)))
	// save state of this actor
	err := impl.SaveState()
	if err != nil {
		return nil, actorErr.ErrSaveStateFailed
	}
	methodType, err := getAbsctractMethodMap(impl)
	if err != nil {
		log.Printf("failed to get absctract method map from registered provider, err = %s", err)
		return nil, actorErr.ErrActorServerInvalid
	}
	return &DefaultActorContainer{
		methodType: methodType,
		actor:      impl,
		serializer: serializer,
	}, actorErr.Success
}

func (d *DefaultActorContainer) GetActor() actor.Server {
	return d.actor
}

// Invoke call actor method with given methodName and param.
func (d *DefaultActorContainer) Invoke(methodName string, param []byte) ([]reflect.Value, actorErr.ActorErr) {
	methodType, ok := d.methodType[methodName]
	if !ok {
		return nil, actorErr.ErrActorMethodNoFound
	}
	argsValues := make([]reflect.Value, 0)
	argsValues = append(argsValues, reflect.ValueOf(d.actor), reflect.ValueOf(context.Background()))
	if len(methodType.argsType) > 0 {
		typ := methodType.argsType[0]
		paramValue := reflect.New(typ)
		paramInterface := paramValue.Interface()
		if err := d.serializer.Unmarshal(param, paramInterface); err != nil {
			return nil, actorErr.ErrActorMethodSerializeFailed
		}
		argsValues = append(argsValues, reflect.ValueOf(paramInterface).Elem())
	}
	returnValue := methodType.method.Func.Call(argsValues)
	return returnValue, actorErr.Success
}
