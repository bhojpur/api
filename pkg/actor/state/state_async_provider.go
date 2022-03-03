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
	"context"

	"github.com/pkg/errors"

	"github.com/bhojpur/api/pkg/actor/codec"
	"github.com/bhojpur/api/pkg/actor/codec/constant"
	client "github.com/bhojpur/api/pkg/client"
)

type AppStateAsyncProvider struct {
	appClient       client.Client
	stateSerializer codec.Codec
}

func (d *AppStateAsyncProvider) Contains(actorType string, actorID string, stateName string) (bool, error) {
	result, err := d.appClient.GetActorState(context.Background(), &client.GetActorStateRequest{
		ActorType: actorType,
		ActorID:   actorID,
		KeyName:   stateName,
	})
	if err != nil || result == nil {
		return false, err
	}
	return len(result.Data) > 0, err
}

func (d *AppStateAsyncProvider) Load(actorType, actorID, stateName string, reply interface{}) error {
	result, err := d.appClient.GetActorState(context.Background(), &client.GetActorStateRequest{
		ActorType: actorType,
		ActorID:   actorID,
		KeyName:   stateName,
	})
	if err != nil {
		return errors.Errorf("get actor state error = %s", err.Error())
	}
	if len(result.Data) == 0 {
		return errors.Errorf("get actor state result empty, with actorType: %s, actorID: %s, stateName %s", actorType, actorID, stateName)
	}
	if err := d.stateSerializer.Unmarshal(result.Data, reply); err != nil {
		return errors.Errorf("unmarshal state data error = %s", err.Error())
	}
	return nil
}

func (d *AppStateAsyncProvider) Apply(actorType, actorID string, changes []*ActorStateChange) error {
	if len(changes) == 0 {
		return nil
	}

	operations := make([]*client.ActorStateOperation, 0)
	var value []byte
	for _, stateChange := range changes {
		if stateChange == nil {
			continue
		}

		appOperationName := string(stateChange.changeKind)
		if len(appOperationName) == 0 {
			continue
		}

		if stateChange.changeKind == Add {
			data, err := d.stateSerializer.Marshal(stateChange.value)
			if err != nil {
				return err
			}
			value = data
		}
		operations = append(operations, &client.ActorStateOperation{
			OperationType: appOperationName,
			Key:           stateChange.stateName,
			Value:         value,
		})
	}

	if len(operations) == 0 {
		return nil
	}

	return d.appClient.SaveStateTransactionally(context.Background(), actorType, actorID, operations)
}

// TODO: the appClient may be nil.
func NewAppStateAsyncProvider(appClient client.Client) *AppStateAsyncProvider {
	stateSerializer, _ := codec.GetActorCodec(constant.DefaultSerializerType)
	return &AppStateAsyncProvider{
		stateSerializer: stateSerializer,
		appClient:       appClient,
	}
}
