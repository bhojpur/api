package internal_test

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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/api/pkg/service/common"
	"github.com/bhojpur/api/pkg/service/internal"
)

func TestTopicRegistrarValidation(t *testing.T) {
	fn := func(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
		return false, nil
	}
	tests := map[string]struct {
		sub common.Subscription
		fn  common.TopicEventHandler
		err string
	}{
		"pubsub required": {
			common.Subscription{ //nolint:exhaustivestruct
				PubsubName: "",
				Topic:      "test",
			}, fn, "pub/sub name required",
		},
		"topic required": {
			common.Subscription{ //nolint:exhaustivestruct
				PubsubName: "test",
				Topic:      "",
			}, fn, "topic name required",
		},
		"handler required": {
			common.Subscription{ //nolint:exhaustivestruct
				PubsubName: "test",
				Topic:      "test",
			}, nil, "topic handler required",
		},
		"route required for routing rule": {
			common.Subscription{ //nolint:exhaustivestruct
				PubsubName: "test",
				Topic:      "test",
				Route:      "",
				Match:      `event.type == "test"`,
			}, fn, "path is required for routing rules",
		},
		"success default route": {
			common.Subscription{ //nolint:exhaustivestruct
				PubsubName: "test",
				Topic:      "test",
			}, fn, "",
		},
		"success routing rule": {
			common.Subscription{ //nolint:exhaustivestruct
				PubsubName: "test",
				Topic:      "test",
				Route:      "/test",
				Match:      `event.type == "test"`,
			}, fn, "",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			m := internal.TopicRegistrar{}
			if tt.err != "" {
				assert.EqualError(t, m.AddSubscription(&tt.sub, tt.fn), tt.err)
			} else {
				assert.NoError(t, m.AddSubscription(&tt.sub, tt.fn))
			}
		})
	}
}

func TestTopicAddSubscriptionMetadata(t *testing.T) {
	handler := func(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
		return false, nil
	}
	topicRegistrar := internal.TopicRegistrar{}
	sub := &common.Subscription{
		PubsubName: "pubsubname",
		Topic:      "topic",
		Metadata:   map[string]string{"key": "value"},
	}

	assert.NoError(t, topicRegistrar.AddSubscription(sub, handler))

	actual := topicRegistrar["pubsubname-topic"].Subscription
	expected := &internal.TopicSubscription{
		PubsubName: sub.PubsubName,
		Topic:      sub.Topic,
		Metadata:   sub.Metadata,
	}
	assert.Equal(t, expected, actual)
}
