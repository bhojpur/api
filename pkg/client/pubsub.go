package client

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
	"encoding/json"
	"log"

	"github.com/pkg/errors"

	pb "github.com/bhojpur/application/pkg/api/v1/runtime"
)

// PublishEventOption is the type for the functional option.
type PublishEventOption func(*pb.PublishEventRequest)

// PublishEvent publishes data onto specific pubsub topic.
func (c *GRPCClient) PublishEvent(ctx context.Context, pubsubName, topicName string, data interface{}, opts ...PublishEventOption) error {
	if pubsubName == "" {
		return errors.New("pubsubName name required")
	}
	if topicName == "" {
		return errors.New("topic name required")
	}

	request := &pb.PublishEventRequest{
		PubsubName: pubsubName,
		Topic:      topicName,
	}
	for _, opt := range opts {
		opt(request)
	}

	if data != nil {
		switch d := data.(type) {
		case []byte:
			request.Data = d
		case string:
			request.Data = []byte(d)
		default:
			var err error
			request.DataContentType = "application/json"
			request.Data, err = json.Marshal(d)
			if err != nil {
				return errors.WithMessage(err, "error serializing input struct")
			}
		}
	}

	_, err := c.protoClient.PublishEvent(c.withAuthToken(ctx), request)
	if err != nil {
		return errors.Wrapf(err, "error publishing event unto %s topic", topicName)
	}

	return nil
}

// PublishEventWithContentType can be passed as option to PublishEvent to set an explicit Content-Type.
func PublishEventWithContentType(contentType string) PublishEventOption {
	return func(e *pb.PublishEventRequest) {
		e.DataContentType = contentType
	}
}

// PublishEventWithMetadata can be passed as option to PublishEvent to set metadata.
func PublishEventWithMetadata(metadata map[string]string) PublishEventOption {
	return func(e *pb.PublishEventRequest) {
		e.Metadata = metadata
	}
}

// PublishEventfromCustomContent serializes an struct and publishes its contents as data (JSON) onto topic in specific pubsub component.
// Deprecated: This method is deprecated and will be removed in a future version of the SDK. Please use `PublishEvent` instead.
func (c *GRPCClient) PublishEventfromCustomContent(ctx context.Context, pubsubName, topicName string, data interface{}) error {
	log.Println("DEPRECATED: client.PublishEventfromCustomContent is deprecated and will be removed in a future version of the SDK. Please use `PublishEvent` instead.")

	// Perform the JSON marshaling here just in case someone passed a []byte or string as data
	enc, err := json.Marshal(data)
	if err != nil {
		return errors.WithMessage(err, "error serializing input struct")
	}

	return c.PublishEvent(ctx, pubsubName, topicName, enc, PublishEventWithContentType("application/json"))
}
