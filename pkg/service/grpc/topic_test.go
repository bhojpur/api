package grpc

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
	"errors"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/api/pkg/service/common"
	runtime "github.com/bhojpur/application/pkg/api/v1/runtime"
)

func TestTopicErrors(t *testing.T) {
	server := getTestServer()
	err := server.AddTopicEventHandler(nil, nil)
	assert.Errorf(t, err, "expected error on nil sub")

	sub := &common.Subscription{}
	err = server.AddTopicEventHandler(sub, nil)
	assert.Errorf(t, err, "expected error on invalid sub")

	sub.PubsubName = "messages"
	err = server.AddTopicEventHandler(sub, nil)
	assert.Errorf(t, err, "expected error on sub without topic")

	sub.Topic = "test"
	err = server.AddTopicEventHandler(sub, nil)
	assert.Errorf(t, err, "expected error on sub without handler")
}

func TestTopicSubscriptionList(t *testing.T) {
	server := getTestServer()

	// Add default route.
	sub1 := &common.Subscription{
		PubsubName: "messages",
		Topic:      "test",
		Route:      "/test",
	}
	err := server.AddTopicEventHandler(sub1, eventHandler)
	assert.Nil(t, err)
	resp, err := server.ListTopicSubscriptions(context.Background(), &empty.Empty{})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	if assert.Lenf(t, resp.Subscriptions, 1, "expected 1 handlers") {
		sub := resp.Subscriptions[0]
		assert.Equal(t, "messages", sub.PubsubName)
		assert.Equal(t, "test", sub.Topic)
		assert.Nil(t, sub.Routes)
	}

	// Add routing rule.
	sub2 := &common.Subscription{
		PubsubName: "messages",
		Topic:      "test",
		Route:      "/other",
		Match:      `event.type == "other"`,
	}
	err = server.AddTopicEventHandler(sub2, eventHandler)
	assert.Nil(t, err)
	resp, err = server.ListTopicSubscriptions(context.Background(), &empty.Empty{})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	if assert.Lenf(t, resp.Subscriptions, 1, "expected 1 handlers") {
		sub := resp.Subscriptions[0]
		assert.Equal(t, "messages", sub.PubsubName)
		assert.Equal(t, "test", sub.Topic)
		if assert.NotNil(t, sub.Routes) {
			assert.Equal(t, "/test", sub.Routes.Default)
			if assert.Len(t, sub.Routes.Rules, 1) {
				rule := sub.Routes.Rules[0]
				assert.Equal(t, "/other", rule.Path)
				assert.Equal(t, `event.type == "other"`, rule.Match)
			}
		}
	}
}

// go test -timeout 30s ./service/grpc -count 1 -run ^TestTopic$
func TestTopic(t *testing.T) {
	ctx := context.Background()

	sub := &common.Subscription{
		PubsubName: "messages",
		Topic:      "test",
	}
	server := getTestServer()

	err := server.AddTopicEventHandler(sub, eventHandler)
	assert.Nil(t, err)

	startTestServer(server)

	t.Run("topic event without request", func(t *testing.T) {
		_, err := server.OnTopicEvent(ctx, nil)
		assert.Error(t, err)
	})

	t.Run("topic event for wrong topic", func(t *testing.T) {
		in := &runtime.TopicEventRequest{
			Topic: "invalid",
		}
		_, err := server.OnTopicEvent(ctx, in)
		assert.Error(t, err)
	})

	t.Run("topic event for valid topic", func(t *testing.T) {
		in := &runtime.TopicEventRequest{
			Id:              "a123",
			Source:          "test",
			Type:            "test",
			SpecVersion:     "v1.0",
			DataContentType: "text/plain",
			Data:            []byte("test"),
			Topic:           sub.Topic,
			PubsubName:      sub.PubsubName,
		}
		_, err := server.OnTopicEvent(ctx, in)
		assert.NoError(t, err)
	})

	stopTestServer(t, server)
}

func TestTopicWithErrors(t *testing.T) {
	ctx := context.Background()

	sub1 := &common.Subscription{
		PubsubName: "messages",
		Topic:      "test1",
	}

	sub2 := &common.Subscription{
		PubsubName: "messages",
		Topic:      "test2",
	}
	server := getTestServer()

	err := server.AddTopicEventHandler(sub1, eventHandlerWithRetryError)
	assert.Nil(t, err)

	err = server.AddTopicEventHandler(sub2, eventHandlerWithError)
	assert.Nil(t, err)

	startTestServer(server)

	t.Run("topic event for retry error", func(t *testing.T) {
		in := &runtime.TopicEventRequest{
			Id:              "a123",
			Source:          "test",
			Type:            "test",
			SpecVersion:     "v1.0",
			DataContentType: "text/plain",
			Data:            []byte("test"),
			Topic:           sub1.Topic,
			PubsubName:      sub1.PubsubName,
		}
		resp, err := server.OnTopicEvent(ctx, in)
		assert.Error(t, err)
		assert.Equal(t, resp.GetStatus(), runtime.TopicEventResponse_RETRY)
	})

	t.Run("topic event for error", func(t *testing.T) {
		in := &runtime.TopicEventRequest{
			Id:              "a123",
			Source:          "test",
			Type:            "test",
			SpecVersion:     "v1.0",
			DataContentType: "text/plain",
			Data:            []byte("test"),
			Topic:           sub2.Topic,
			PubsubName:      sub2.PubsubName,
		}
		resp, err := server.OnTopicEvent(ctx, in)
		assert.Error(t, err)
		assert.Equal(t, resp.GetStatus(), runtime.TopicEventResponse_DROP)
	})

	stopTestServer(t, server)
}

func eventHandler(ctx context.Context, event *common.TopicEvent) (retry bool, err error) {
	if event == nil {
		return true, errors.New("nil event")
	}
	return false, nil
}

func eventHandlerWithRetryError(ctx context.Context, event *common.TopicEvent) (retry bool, err error) {
	return true, errors.New("nil event")
}

func eventHandlerWithError(ctx context.Context, event *common.TopicEvent) (retry bool, err error) {
	return false, errors.New("nil event")
}

func TestEventDataHandling(t *testing.T) {
	ctx := context.Background()

	tests := map[string]struct {
		contentType string
		data        string
		value       interface{}
	}{
		"JSON bytes": {
			contentType: "application/json",
			data:        `{"message":"hello"}`,
			value: map[string]interface{}{
				"message": "hello",
			},
		},
		"JSON entension media type bytes": {
			contentType: "application/extension+json",
			data:        `{"message":"hello"}`,
			value: map[string]interface{}{
				"message": "hello",
			},
		},
		"Test": {
			contentType: "text/plain",
			data:        `message = hello`,
			value:       `message = hello`,
		},
		"Other": {
			contentType: "application/octet-stream",
			data:        `message = hello`,
			value:       []byte(`message = hello`),
		},
	}

	s := getTestServer()

	sub := &common.Subscription{
		PubsubName: "messages",
		Topic:      "test",
		Route:      "/test",
		Metadata:   map[string]string{},
	}

	recv := make(chan struct{}, 1)
	var topicEvent *common.TopicEvent
	handler := func(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
		topicEvent = e
		recv <- struct{}{}

		return false, nil
	}
	err := s.AddTopicEventHandler(sub, handler)
	assert.NoErrorf(t, err, "error adding event handler")

	startTestServer(s)

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			in := runtime.TopicEventRequest{
				Id:              "a123",
				Source:          "test",
				Type:            "test",
				SpecVersion:     "v1.0",
				DataContentType: tt.contentType,
				Data:            []byte(tt.data),
				Topic:           sub.Topic,
				PubsubName:      sub.PubsubName,
			}

			s.OnTopicEvent(ctx, &in)
			<-recv
			assert.Equal(t, tt.value, topicEvent.Data)
		})
	}
}
