package main

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
	"fmt"
	"log"
	"net/http"

	"github.com/bhojpur/api/internal/basic/protocol"
	"github.com/bhojpur/api/pkg/actor"
	app "github.com/bhojpur/api/pkg/client"

	appsvc "github.com/bhojpur/api/pkg/service/http"
)

func testActorFactory() actor.Server {
	client, err := app.NewClient()
	if err != nil {
		panic(err)
	}
	return &TestActor{
		appClient: client,
	}
}

type TestActor struct {
	actor.ServerImplBase
	appClient app.Client
}

func (t *TestActor) Type() string {
	return "testActorType"
}

// user defined functions
func (t *TestActor) StopTimer(ctx context.Context, req *protocol.TimerRequest) error {
	return t.appClient.UnregisterActorTimer(ctx, &app.UnregisterActorTimerRequest{
		ActorType: t.Type(),
		ActorID:   t.ID(),
		Name:      req.TimerName,
	})
}

func (t *TestActor) StartTimer(ctx context.Context, req *protocol.TimerRequest) error {
	return t.appClient.RegisterActorTimer(ctx, &app.RegisterActorTimerRequest{
		ActorType: t.Type(),
		ActorID:   t.ID(),
		Name:      req.TimerName,
		DueTime:   req.Duration,
		Period:    req.Period,
		Data:      []byte(req.Data),
		CallBack:  req.CallBack,
	})
}

func (t *TestActor) StartReminder(ctx context.Context, req *protocol.ReminderRequest) error {
	return t.appClient.RegisterActorReminder(ctx, &app.RegisterActorReminderRequest{
		ActorType: t.Type(),
		ActorID:   t.ID(),
		Name:      req.ReminderName,
		DueTime:   req.Duration,
		Period:    req.Period,
		Data:      []byte(req.Data),
	})
}

func (t *TestActor) StopReminder(ctx context.Context, req *protocol.ReminderRequest) error {
	return t.appClient.UnregisterActorReminder(ctx, &app.UnregisterActorReminderRequest{
		ActorType: t.Type(),
		ActorID:   t.ID(),
		Name:      req.ReminderName,
	})
}

func (t *TestActor) Invoke(ctx context.Context, req string) (string, error) {
	fmt.Println("get req = ", req)
	return req, nil
}

func (t *TestActor) GetUser(ctx context.Context, user *protocol.User) (*protocol.User, error) {
	fmt.Println("call get user req = ", user)
	return user, nil
}

func (t *TestActor) Get(context.Context) (string, error) {
	return "get result", nil
}

func (t *TestActor) Post(ctx context.Context, req string) error {
	fmt.Println("get post request = ", req)
	return nil
}

func (t *TestActor) IncrementAndGet(ctx context.Context, stateKey string) (*protocol.User, error) {
	stateData := protocol.User{}
	if exist, err := t.GetStateManager().Contains(stateKey); err != nil {
		fmt.Println("state manager call contains with key " + stateKey + "err = " + err.Error())
		return &stateData, err
	} else if exist {
		if err := t.GetStateManager().Get(stateKey, &stateData); err != nil {
			fmt.Println("state manager call get with key " + stateKey + "err = " + err.Error())
			return &stateData, err
		}
	}
	stateData.Age++
	if err := t.GetStateManager().Set(stateKey, stateData); err != nil {
		fmt.Printf("state manager set get with key %s and state data = %+v, error = %s", stateKey, stateData, err.Error())
		return &stateData, err
	}
	return &stateData, nil
}

func (t *TestActor) ReminderCall(reminderName string, state []byte, dueTime string, period string) {
	fmt.Println("receive reminder = ", reminderName, " state = ", string(state), "duetime = ", dueTime, "period = ", period)
}

func main() {
	s := appsvc.NewService(":3000")
	s.RegisterActorImplFactory(testActorFactory)
	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error listenning: %v", err)
	}
}
