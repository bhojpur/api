package state

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
	"reflect"
	"sync"

	"github.com/pkg/errors"

	"github.com/bhojpur/api/pkg/actor"
)

type ActorStateManager struct {
	ActorTypeName      string
	ActorID            string
	stateChangeTracker sync.Map // map[string]*ChangeMetadata
	stateAsyncProvider *AppStateAsyncProvider
}

func (a *ActorStateManager) Add(stateName string, value interface{}) error {
	if stateName == "" {
		return errors.Errorf("state's name can't be empty")
	}
	exists, err := a.stateAsyncProvider.Contains(a.ActorTypeName, a.ActorID, stateName)
	if err != nil {
		return err
	}

	if val, ok := a.stateChangeTracker.Load(stateName); ok {
		metadata := val.(*ChangeMetadata)
		if metadata.Kind == Remove {
			a.stateChangeTracker.Store(stateName, &ChangeMetadata{
				Kind:  Update,
				Value: value,
			})
			return nil
		}
		return errors.Errorf("Duplicate cached state: %s", stateName)
	}
	if exists {
		return errors.Errorf("Duplicate state: %s", stateName)
	}
	a.stateChangeTracker.Store(stateName, &ChangeMetadata{
		Kind:  Add,
		Value: value,
	})
	return nil
}

func (a *ActorStateManager) Get(stateName string, reply interface{}) error {
	if stateName == "" {
		return errors.Errorf("state's name can't be empty")
	}

	if val, ok := a.stateChangeTracker.Load(stateName); ok {
		metadata := val.(*ChangeMetadata)
		if metadata.Kind == Remove {
			return errors.Errorf("state is marked for remove: %s", stateName)
		}
		replyVal := reflect.ValueOf(reply).Elem()
		metadataValue := reflect.ValueOf(metadata.Value)
		if metadataValue.Kind() == reflect.Ptr {
			replyVal.Set(metadataValue.Elem())
		} else {
			replyVal.Set(metadataValue)
		}

		return nil
	}

	err := a.stateAsyncProvider.Load(a.ActorTypeName, a.ActorID, stateName, reply)
	a.stateChangeTracker.Store(stateName, &ChangeMetadata{
		Kind:  None,
		Value: reply,
	})
	return err
}

func (a *ActorStateManager) Set(stateName string, value interface{}) error {
	if stateName == "" {
		return errors.Errorf("state's name can't be empty")
	}
	if val, ok := a.stateChangeTracker.Load(stateName); ok {
		metadata := val.(*ChangeMetadata)
		if metadata.Kind == None || metadata.Kind == Remove {
			metadata.Kind = Update
		}
		a.stateChangeTracker.Store(stateName, NewChangeMetadata(metadata.Kind, value))
		return nil
	}
	a.stateChangeTracker.Store(stateName, &ChangeMetadata{
		Kind:  Add,
		Value: value,
	})
	return nil
}

func (a *ActorStateManager) Remove(stateName string) error {
	if stateName == "" {
		return errors.Errorf("state's name can't be empty")
	}
	if val, ok := a.stateChangeTracker.Load(stateName); ok {
		metadata := val.(*ChangeMetadata)
		if metadata.Kind == Remove {
			return nil
		}
		if metadata.Kind == Add {
			a.stateChangeTracker.Delete(stateName)
			return nil
		}

		a.stateChangeTracker.Store(stateName, &ChangeMetadata{
			Kind:  Remove,
			Value: nil,
		})
		return nil
	}
	if exist, err := a.stateAsyncProvider.Contains(a.ActorTypeName, a.ActorID, stateName); err != nil && exist {
		a.stateChangeTracker.Store(stateName, &ChangeMetadata{
			Kind:  Remove,
			Value: nil,
		})
	}
	return nil
}

func (a *ActorStateManager) Contains(stateName string) (bool, error) {
	if stateName == "" {
		return false, errors.Errorf("state's name can't be empty")
	}
	if val, ok := a.stateChangeTracker.Load(stateName); ok {
		metadata := val.(*ChangeMetadata)
		if metadata.Kind == Remove {
			return false, nil
		}
		return true, nil
	}
	return a.stateAsyncProvider.Contains(a.ActorTypeName, a.ActorID, stateName)
}

func (a *ActorStateManager) Save() error {
	changes := make([]*ActorStateChange, 0)
	a.stateChangeTracker.Range(func(key, value interface{}) bool {
		stateName := key.(string)
		metadata := value.(*ChangeMetadata)
		changes = append(changes, NewActorStateChange(stateName, metadata.Value, metadata.Kind))
		return true
	})
	if err := a.stateAsyncProvider.Apply(a.ActorTypeName, a.ActorID, changes); err != nil {
		return err
	}
	a.Flush()
	return nil
}

func (a *ActorStateManager) Flush() {
	a.stateChangeTracker.Range(func(key, value interface{}) bool {
		stateName := key.(string)
		metadata := value.(*ChangeMetadata)
		if metadata.Kind == Remove {
			a.stateChangeTracker.Delete(stateName)
			return true
		}
		metadata = NewChangeMetadata(None, metadata.Value)
		a.stateChangeTracker.Store(stateName, metadata)
		return true
	})
}

func NewActorStateManager(actorTypeName string, actorID string, provider *AppStateAsyncProvider) actor.StateManager {
	return &ActorStateManager{
		stateAsyncProvider: provider,
		ActorTypeName:      actorTypeName,
		ActorID:            actorID,
	}
}
