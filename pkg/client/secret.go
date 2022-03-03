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

	"github.com/pkg/errors"

	pb "github.com/bhojpur/application/pkg/api/v1/runtime"
)

// GetSecret retrieves preconfigured secret from specified store using key.
func (c *GRPCClient) GetSecret(ctx context.Context, storeName, key string, meta map[string]string) (data map[string]string, err error) {
	if storeName == "" {
		return nil, errors.New("nil storeName")
	}
	if key == "" {
		return nil, errors.New("nil key")
	}

	req := &pb.GetSecretRequest{
		Key:       key,
		StoreName: storeName,
		Metadata:  meta,
	}

	resp, err := c.protoClient.GetSecret(c.withAuthToken(ctx), req)
	if err != nil {
		return nil, errors.Wrap(err, "error invoking service")
	}

	if resp != nil {
		data = resp.GetData()
	}

	return
}

// GetBulkSecret retrieves all preconfigured secrets for this application.
func (c *GRPCClient) GetBulkSecret(ctx context.Context, storeName string, meta map[string]string) (data map[string]map[string]string, err error) {
	if storeName == "" {
		return nil, errors.New("nil storeName")
	}

	req := &pb.GetBulkSecretRequest{
		StoreName: storeName,
		Metadata:  meta,
	}

	resp, err := c.protoClient.GetBulkSecret(c.withAuthToken(ctx), req)
	if err != nil {
		return nil, errors.Wrap(err, "error invoking service")
	}

	if resp != nil {
		data = map[string]map[string]string{}

		for secretName, secretResponse := range resp.Data {
			data[secretName] = map[string]string{}

			for k, v := range secretResponse.Secrets {
				data[secretName][k] = v
			}
		}
	}

	return
}
