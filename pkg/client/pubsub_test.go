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
	"testing"

	"github.com/stretchr/testify/assert"
)

type _testCustomContentwithText struct {
	Key1, Key2 string
}

type _testCustomContentwithTextandNumbers struct {
	Key1 string
	Key2 int
}

type _testCustomContentwithSlices struct {
	Key1 []string
	Key2 []int
}

// go test -timeout 30s ./client -count 1 -run ^TestPublishEvent$
func TestPublishEvent(t *testing.T) {
	ctx := context.Background()

	t.Run("with data", func(t *testing.T) {
		err := testClient.PublishEvent(ctx, "messages", "test", []byte("ping"))
		assert.Nil(t, err)
	})

	t.Run("without data", func(t *testing.T) {
		err := testClient.PublishEvent(ctx, "messages", "test", nil)
		assert.Nil(t, err)
	})

	t.Run("with empty topic name", func(t *testing.T) {
		err := testClient.PublishEvent(ctx, "messages", "", []byte("ping"))
		assert.NotNil(t, err)
	})

	t.Run("from struct with text", func(t *testing.T) {
		testdata := _testStructwithText{
			Key1: "value1",
			Key2: "value2",
		}
		err := testClient.PublishEventfromCustomContent(ctx, "messages", "test", testdata)
		assert.Nil(t, err)
	})

	t.Run("from struct with text and numbers", func(t *testing.T) {
		testdata := _testStructwithTextandNumbers{
			Key1: "value1",
			Key2: 2500,
		}
		err := testClient.PublishEventfromCustomContent(ctx, "messages", "test", testdata)
		assert.Nil(t, err)
	})

	t.Run("from struct with slices", func(t *testing.T) {
		testdata := _testStructwithSlices{
			Key1: []string{"value1", "value2", "value3"},
			Key2: []int{25, 40, 600},
		}
		err := testClient.PublishEventfromCustomContent(ctx, "messages", "test", testdata)
		assert.Nil(t, err)
	})

	t.Run("error serializing JSON", func(t *testing.T) {
		err := testClient.PublishEventfromCustomContent(ctx, "messages", "test", make(chan struct{}))
		assert.Error(t, err)
	})
}
