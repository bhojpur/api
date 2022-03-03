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
	"testing"

	"github.com/bhojpur/api/pkg/actor/codec"
	"github.com/bhojpur/api/pkg/client"
)

func TestAppStateAsyncProvider_Apply(t *testing.T) {
	type fields struct {
		appClient       client.Client
		stateSerializer codec.Codec
	}
	type args struct {
		actorType string
		actorID   string
		changes   []*ActorStateChange
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "empty changes",
			args: args{
				actorType: "testActor",
				actorID:   "test-0",
				changes:   nil,
			},
			wantErr: false,
		},
		{
			name: "only readonly state changes",
			args: args{
				actorType: "testActor",
				actorID:   "test-0",
				changes: []*ActorStateChange{
					{
						stateName:  "stateName1",
						value:      "Any",
						changeKind: None,
					},
					{
						stateName:  "stateName2",
						value:      "Any",
						changeKind: None,
					},
				},
			},
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &AppStateAsyncProvider{
				appClient:       tt.fields.appClient,
				stateSerializer: tt.fields.stateSerializer,
			}
			if err := d.Apply(tt.args.actorType, tt.args.actorID, tt.args.changes); (err != nil) != tt.wantErr {
				t.Errorf("Apply() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAppStateAsyncProvider_Contains(t *testing.T) {
	type fields struct {
		appClient       client.Client
		stateSerializer codec.Codec
	}
	type args struct {
		actorType string
		actorID   string
		stateName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &AppStateAsyncProvider{
				appClient:       tt.fields.appClient,
				stateSerializer: tt.fields.stateSerializer,
			}
			got, err := d.Contains(tt.args.actorType, tt.args.actorID, tt.args.stateName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Contains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Contains() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppStateAsyncProvider_Load(t *testing.T) {
	type fields struct {
		appClient       client.Client
		stateSerializer codec.Codec
	}
	type args struct {
		actorType string
		actorID   string
		stateName string
		reply     interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &AppStateAsyncProvider{
				appClient:       tt.fields.appClient,
				stateSerializer: tt.fields.stateSerializer,
			}
			if err := d.Load(tt.args.actorType, tt.args.actorID, tt.args.stateName, tt.args.reply); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewAppStateAsyncProvider(t *testing.T) {
	type args struct {
		appClient client.Client
	}
	tests := []struct {
		name string
		args args
		want *AppStateAsyncProvider
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAppStateAsyncProvider(tt.args.appClient); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAppStateAsyncProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}
