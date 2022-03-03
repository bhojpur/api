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
	"time"

	"github.com/bhojpur/api/internal/basic/protocol"
	app "github.com/bhojpur/api/pkg/client"
)

func main() {
	ctx := context.Background()

	// create the client
	client, err := app.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// implement actor client stub
	myActor := new(protocol.ClientStub)
	client.ImplActorClientStub(myActor)

	// Invoke user defined method GetUser with user defined param protocol.User and response
	// using default serializer type json
	user, err := myActor.GetUser(ctx, &protocol.User{
		Name: "abc",
		Age:  123,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("get user result = ", user)

	// Invoke user defined method Invoke
	rsp, err := myActor.Invoke(ctx, "pramila")
	if err != nil {
		panic(err)
	}
	fmt.Println("get invoke result = ", rsp)

	// Invoke user defined method Post with empty response
	err = myActor.Post(ctx, "pramila")
	if err != nil {
		panic(err)
	}
	fmt.Println("get post result = ", rsp)

	// Invoke user defined method Get with empty request param
	rsp, err = myActor.Get(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("get result = ", rsp)

	// Invoke user defined method StarTimer, and server side actor start actor timer with given params.
	err = myActor.StartTimer(ctx, &protocol.TimerRequest{
		TimerName: "testTimerName",
		CallBack:  "Invoke",
		Period:    "5s",
		Duration:  "5s",
		Data:      `"hello"`,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("start timer")
	<-time.After(time.Second * 10) // timer call for twice

	// Invoke user defined method StopTimer, and server side actor stop actor timer with given params.
	err = myActor.StopTimer(ctx, &protocol.TimerRequest{
		TimerName: "testTimerName",
		CallBack:  "Invoke",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("stop timer")

	// Invoke user defined method StartReminder, and server side actor start actor reminder with given params.
	err = myActor.StartReminder(ctx, &protocol.ReminderRequest{
		ReminderName: "testReminderName",
		Period:       "5s",
		Duration:     "5s",
		Data:         `"hello"`,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("start reminder")
	<-time.After(time.Second * 10) // timer call for twice

	// Invoke user defined method StopReminder, and server side actor stop actor reminder with given params.
	err = myActor.StopReminder(ctx, &protocol.ReminderRequest{
		ReminderName: "testReminderName",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("stop reminder")

	for i := 0; i < 2; i++ {
		// Invoke user defined method IncrementAndGet, and server side actor increase the state named testStateKey and return.
		usr, err := myActor.IncrementAndGet(ctx, "testStateKey")
		if err != nil {
			panic(err)
		}
		fmt.Printf("get user = %+v\n", *usr)
		time.Sleep(time.Second)
	}
}
