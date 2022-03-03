package internal

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
	"errors"
	"fmt"

	"github.com/bhojpur/api/pkg/service/common"
)

// TopicRegistrar is a map of <pubsubname>-<topic> to `TopicRegistration`
// and acts as a lookup as the application is building up subscriptions with
// potentially multiple routes per topic.
type TopicRegistrar map[string]*TopicRegistration

// TopicRegistration encapsulates the subscription and handlers.
type TopicRegistration struct {
	Subscription   *TopicSubscription
	DefaultHandler common.TopicEventHandler
	RouteHandlers  map[string]common.TopicEventHandler
}

func (m TopicRegistrar) AddSubscription(sub *common.Subscription, fn common.TopicEventHandler) error {
	if sub.Topic == "" {
		return errors.New("topic name required")
	}
	if sub.PubsubName == "" {
		return errors.New("pub/sub name required")
	}
	if fn == nil {
		return fmt.Errorf("topic handler required")
	}
	key := sub.PubsubName + "-" + sub.Topic
	ts, ok := m[key]
	if !ok {
		ts = &TopicRegistration{
			Subscription:   NewTopicSubscription(sub.PubsubName, sub.Topic),
			RouteHandlers:  make(map[string]common.TopicEventHandler),
			DefaultHandler: nil,
		}
		ts.Subscription.SetMetadata(sub.Metadata)
		m[key] = ts
	}

	if sub.Match != "" {
		if err := ts.Subscription.AddRoutingRule(sub.Route, sub.Match, sub.Priority); err != nil {
			return err
		}
	} else {
		if err := ts.Subscription.SetDefaultRoute(sub.Route); err != nil {
			return err
		}
		ts.DefaultHandler = fn
	}
	ts.RouteHandlers[sub.Route] = fn

	return nil
}
