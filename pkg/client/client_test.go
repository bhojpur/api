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
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/anypb"

	commonv1pb "github.com/bhojpur/application/pkg/api/v1/common"
	pb "github.com/bhojpur/application/pkg/api/v1/runtime"
)

const (
	testBufSize = 1024 * 1024
	testSocket  = "/tmp/app.socket"
)

var testClient Client

func TestMain(m *testing.M) {
	ctx := context.Background()
	c, f := getTestClient(ctx)
	testClient = c
	r := m.Run()
	f()

	if r != 0 {
		os.Exit(r)
	}

	c, f = getTestClientWithSocket(ctx)
	testClient = c
	r = m.Run()
	f()
	os.Exit(r)
}

func TestNewClient(t *testing.T) {
	t.Run("return error when unable to reach server", func(t *testing.T) {
		_, err := NewClientWithPort("1")
		assert.Error(t, err)
	})

	t.Run("no arg for with port", func(t *testing.T) {
		_, err := NewClientWithPort("")
		assert.Error(t, err)
	})

	t.Run("no arg for with address", func(t *testing.T) {
		_, err := NewClientWithAddress("")
		assert.Error(t, err)
	})

	t.Run("no arg with socket", func(t *testing.T) {
		_, err := NewClientWithSocket("")
		assert.Error(t, err)
	})

	t.Run("new client closed with token", func(t *testing.T) {
		t.Setenv(apiTokenEnvVarName, "test")
		c, err := NewClientWithSocket(testSocket)
		assert.NoError(t, err)
		defer c.Close()
		c.WithAuthToken("")
	})

	t.Run("new client closed with empty token", func(t *testing.T) {
		testClient.WithAuthToken("")
	})

	t.Run("new client with trace ID", func(t *testing.T) {
		_ = testClient.WithTraceID(context.Background(), "test")
	})

	t.Run("new socket client closed with token", func(t *testing.T) {
		t.Setenv(apiTokenEnvVarName, "test")
		c, err := NewClientWithSocket(testSocket)
		assert.NoError(t, err)
		defer c.Close()
		c.WithAuthToken("")
	})

	t.Run("new socket client closed with empty token", func(t *testing.T) {
		c, err := NewClientWithSocket(testSocket)
		assert.NoError(t, err)
		defer c.Close()
		c.WithAuthToken("")
	})

	t.Run("new socket client with trace ID", func(t *testing.T) {
		c, err := NewClientWithSocket(testSocket)
		assert.NoError(t, err)
		defer c.Close()
		ctx := c.WithTraceID(context.Background(), "")
		_ = c.WithTraceID(ctx, "test")
	})
}

func TestShutdown(t *testing.T) {
	ctx := context.Background()

	t.Run("shutdown", func(t *testing.T) {
		err := testClient.Shutdown(ctx)
		assert.NoError(t, err)
	})
}

func getTestClient(ctx context.Context) (client Client, closer func()) {
	s := grpc.NewServer()
	pb.RegisterApplicationServer(s, &testAppServer{
		state: make(map[string][]byte),
	})

	l := bufconn.Listen(testBufSize)
	go func() {
		if err := s.Serve(l); err != nil && err.Error() != "closed" {
			logger.Fatalf("test server exited with error: %v", err)
		}
	}()

	d := grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return l.Dial()
	})

	c, err := grpc.DialContext(ctx, "", d, grpc.WithInsecure())
	if err != nil {
		logger.Fatalf("failed to dial test context: %v", err)
	}

	closer = func() {
		l.Close()
		s.Stop()
	}

	client = NewClientWithConnection(c)
	return
}

func getTestClientWithSocket(ctx context.Context) (client Client, closer func()) {
	s := grpc.NewServer()
	pb.RegisterApplicationServer(s, &testAppServer{
		state: make(map[string][]byte),
	})

	var lc net.ListenConfig
	l, err := lc.Listen(ctx, "unix", testSocket)
	if err != nil {
		logger.Fatalf("socket test server created with error: %v", err)
	}

	go func() {
		if err = s.Serve(l); err != nil && err.Error() != "accept unix /tmp/app.socket: use of closed network connection" {
			logger.Fatalf("socket test server exited with error: %v", err)
		}
	}()

	closer = func() {
		l.Close()
		s.Stop()
		os.Remove(testSocket)
	}

	if client, err = NewClientWithSocket(testSocket); err != nil {
		logger.Fatalf("socket test client created with error: %v", err)
	}

	return
}

type testAppServer struct {
	pb.UnimplementedApplicationServer
	state map[string][]byte
}

func (s *testAppServer) InvokeService(ctx context.Context, req *pb.InvokeServiceRequest) (*commonv1pb.InvokeResponse, error) {
	if req.Message == nil {
		return &commonv1pb.InvokeResponse{
			ContentType: "text/plain",
			Data: &anypb.Any{
				Value: []byte("pong"),
			},
		}, nil
	}
	return &commonv1pb.InvokeResponse{
		ContentType: req.Message.ContentType,
		Data:        req.Message.Data,
	}, nil
}

func (s *testAppServer) GetState(ctx context.Context, req *pb.GetStateRequest) (*pb.GetStateResponse, error) {
	return &pb.GetStateResponse{
		Data: s.state[req.Key],
		Etag: "1",
	}, nil
}

func (s *testAppServer) GetBulkState(ctx context.Context, in *pb.GetBulkStateRequest) (*pb.GetBulkStateResponse, error) {
	items := make([]*pb.BulkStateItem, 0)
	for _, k := range in.GetKeys() {
		if v, found := s.state[k]; found {
			item := &pb.BulkStateItem{
				Key:  k,
				Etag: "1",
				Data: v,
			}
			items = append(items, item)
		}
	}
	return &pb.GetBulkStateResponse{
		Items: items,
	}, nil
}

func (s *testAppServer) SaveState(ctx context.Context, req *pb.SaveStateRequest) (*empty.Empty, error) {
	for _, item := range req.States {
		s.state[item.Key] = item.Value
	}
	return &empty.Empty{}, nil
}

func (s *testAppServer) QueryStateAlpha1(ctx context.Context, req *pb.QueryStateRequest) (*pb.QueryStateResponse, error) {
	var v map[string]interface{}
	if err := json.Unmarshal([]byte(req.Query), &v); err != nil {
		return nil, err
	}

	ret := &pb.QueryStateResponse{
		Results: make([]*pb.QueryStateItem, 0, len(s.state)),
	}
	for key, value := range s.state {
		ret.Results = append(ret.Results, &pb.QueryStateItem{Key: key, Data: value})
	}
	return ret, nil
}

func (s *testAppServer) DeleteState(ctx context.Context, req *pb.DeleteStateRequest) (*empty.Empty, error) {
	delete(s.state, req.Key)
	return &empty.Empty{}, nil
}

func (s *testAppServer) DeleteBulkState(ctx context.Context, req *pb.DeleteBulkStateRequest) (*empty.Empty, error) {
	for _, item := range req.States {
		delete(s.state, item.Key)
	}
	return &empty.Empty{}, nil
}

func (s *testAppServer) ExecuteStateTransaction(ctx context.Context, in *pb.ExecuteStateTransactionRequest) (*empty.Empty, error) {
	for _, op := range in.GetOperations() {
		item := op.GetRequest()
		switch opType := op.GetOperationType(); opType {
		case "upsert":
			s.state[item.Key] = item.Value
		case "delete":
			delete(s.state, item.Key)
		default:
			return &empty.Empty{}, fmt.Errorf("invalid operation type: %s", opType)
		}
	}
	return &empty.Empty{}, nil
}

func (s *testAppServer) PublishEvent(ctx context.Context, req *pb.PublishEventRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (s *testAppServer) InvokeBinding(ctx context.Context, req *pb.InvokeBindingRequest) (*pb.InvokeBindingResponse, error) {
	if req.Data == nil {
		return &pb.InvokeBindingResponse{
			Data:     []byte("test"),
			Metadata: map[string]string{"k1": "v1", "k2": "v2"},
		}, nil
	}
	return &pb.InvokeBindingResponse{
		Data:     req.Data,
		Metadata: req.Metadata,
	}, nil
}

func (s *testAppServer) GetSecret(ctx context.Context, req *pb.GetSecretRequest) (*pb.GetSecretResponse, error) {
	d := make(map[string]string)
	d["test"] = "value"
	return &pb.GetSecretResponse{
		Data: d,
	}, nil
}

func (s *testAppServer) GetBulkSecret(ctx context.Context, req *pb.GetBulkSecretRequest) (*pb.GetBulkSecretResponse, error) {
	d := make(map[string]*pb.SecretResponse)
	d["test"] = &pb.SecretResponse{
		Secrets: map[string]string{
			"test": "value",
		},
	}
	return &pb.GetBulkSecretResponse{
		Data: d,
	}, nil
}

func (s *testAppServer) RegisterActorReminder(ctx context.Context, req *pb.RegisterActorReminderRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (s *testAppServer) UnregisterActorReminder(ctx context.Context, req *pb.UnregisterActorReminderRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (s *testAppServer) RenameActorReminder(ctx context.Context, req *pb.RenameActorReminderRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (s *testAppServer) InvokeActor(context.Context, *pb.InvokeActorRequest) (*pb.InvokeActorResponse, error) {
	return &pb.InvokeActorResponse{
		Data: []byte("mockValue"),
	}, nil
}

func (s *testAppServer) RegisterActorTimer(context.Context, *pb.RegisterActorTimerRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (s *testAppServer) UnregisterActorTimer(context.Context, *pb.UnregisterActorTimerRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (s *testAppServer) Shutdown(ctx context.Context, req *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}
