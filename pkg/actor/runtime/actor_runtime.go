package runtime

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
	"sync"

	"github.com/bhojpur/api/pkg/actor"
	"github.com/bhojpur/api/pkg/actor/protocol"
	"github.com/bhojpur/api/pkg/actor/config"
	actorErr "github.com/bhojpur/api/pkg/actor/error"
	"github.com/bhojpur/api/pkg/actor/manager"
)

type ActorRunTime struct {
	config        protocol.ActorRuntimeConfig
	actorManagers sync.Map
}

var actorRuntimeInstance *ActorRunTime

// NewActorRuntime creates an empty ActorRuntime.
func NewActorRuntime() *ActorRunTime {
	return &ActorRunTime{}
}

// GetActorRuntimeInstance gets or create runtime instance.
func GetActorRuntimeInstance() *ActorRunTime {
	if actorRuntimeInstance == nil {
		actorRuntimeInstance = NewActorRuntime()
	}
	return actorRuntimeInstance
}

// RegisterActorFactory registers the given actor factory from user, and create new actor manager if not exists.
func (r *ActorRunTime) RegisterActorFactory(f actor.Factory, opt ...config.Option) {
	conf := config.GetConfigFromOptions(opt...)
	actType := f().Type()
	r.config.RegisteredActorTypes = append(r.config.RegisteredActorTypes, actType)
	mng, ok := r.actorManagers.Load(actType)
	if !ok {
		newMng, err := manager.NewDefaultActorManager(conf.SerializerType)
		if err != actorErr.Success {
			return
		}
		newMng.RegisterActorImplFactory(f)
		r.actorManagers.Store(actType, newMng)
		return
	}
	mng.(manager.ActorManager).RegisterActorImplFactory(f)
}

func (r *ActorRunTime) GetJSONSerializedConfig() ([]byte, error) {
	data, err := json.Marshal(&r.config)
	return data, err
}

func (r *ActorRunTime) InvokeActorMethod(actorTypeName, actorID, actorMethod string, payload []byte) ([]byte, actorErr.ActorErr) {
	mng, ok := r.actorManagers.Load(actorTypeName)
	if !ok {
		return nil, actorErr.ErrActorTypeNotFound
	}
	return mng.(manager.ActorManager).InvokeMethod(actorID, actorMethod, payload)
}

func (r *ActorRunTime) Deactivate(actorTypeName, actorID string) actorErr.ActorErr {
	targetManager, ok := r.actorManagers.Load(actorTypeName)
	if !ok {
		return actorErr.ErrActorTypeNotFound
	}
	return targetManager.(manager.ActorManager).DeactivateActor(actorID)
}

func (r *ActorRunTime) InvokeReminder(actorTypeName, actorID, reminderName string, params []byte) actorErr.ActorErr {
	targetManager, ok := r.actorManagers.Load(actorTypeName)
	if !ok {
		return actorErr.ErrActorTypeNotFound
	}
	mng := targetManager.(manager.ActorManager)
	return mng.InvokeReminder(actorID, reminderName, params)
}

func (r *ActorRunTime) InvokeTimer(actorTypeName, actorID, timerName string, params []byte) actorErr.ActorErr {
	targetManager, ok := r.actorManagers.Load(actorTypeName)
	if !ok {
		return actorErr.ErrActorTypeNotFound
	}
	mng := targetManager.(manager.ActorManager)
	return mng.InvokeTimer(actorID, timerName, params)
}
