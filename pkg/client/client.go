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
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/bhojpur/api/pkg/actor"
	"github.com/bhojpur/api/pkg/actor/config"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/bhojpur/application/pkg/api/v1/runtime"

	// used to import codec implements.
	_ "github.com/bhojpur/api/pkg/actor/codec/impl"
)

const (
	appPortDefault     = "50001"
	appPortEnvVarName  = "APP_GRPC_PORT" /* #nosec */
	traceparentKey     = "traceparent"
	apiTokenKey        = "app-api-token" /* #nosec */
	apiTokenEnvVarName = "APP_API_TOKEN" /* #nosec */
)

var (
	logger               = log.New(os.Stdout, "", 0)
	_             Client = (*GRPCClient)(nil)
	defaultClient Client
	doOnce        sync.Once
)

// Client is the interface for Bhojpur Application runtime client implementation.
type Client interface {
	// InvokeBinding invokes specific operation on the configured Bhojpur Application runtime binding.
	// This method covers input, output, and bi-directional bindings.
	InvokeBinding(ctx context.Context, in *InvokeBindingRequest) (out *BindingEvent, err error)

	// InvokeOutputBinding invokes configured Bhojpur Application runtime binding with data.InvokeOutputBinding
	// This method differs from InvokeBinding in that it doesn't expect any content being returned from the invoked method.
	InvokeOutputBinding(ctx context.Context, in *InvokeBindingRequest) error

	// InvokeMethod invokes service without raw data
	InvokeMethod(ctx context.Context, appID, methodName, verb string) (out []byte, err error)

	// InvokeMethodWithContent invokes service with content
	InvokeMethodWithContent(ctx context.Context, appID, methodName, verb string, content *DataContent) (out []byte, err error)

	// InvokeMethodWithCustomContent invokes app with custom content (struct + content type).
	InvokeMethodWithCustomContent(ctx context.Context, appID, methodName, verb string, contentType string, content interface{}) (out []byte, err error)

	// PublishEvent publishes data onto topic in specific pubsub component.
	PublishEvent(ctx context.Context, pubsubName, topicName string, data interface{}, opts ...PublishEventOption) error

	// PublishEventfromCustomContent serializes an struct and publishes its contents as data (JSON) onto topic in specific pubsub component.
	// Deprecated: This method is deprecated and will be removed in a future version of the SDK. Please use `PublishEvent` instead.
	PublishEventfromCustomContent(ctx context.Context, pubsubName, topicName string, data interface{}) error

	// GetSecret retrieves preconfigured secret from specified store using key.
	GetSecret(ctx context.Context, storeName, key string, meta map[string]string) (data map[string]string, err error)

	// GetBulkSecret retrieves all preconfigured secrets for this Bhojpur Application.
	GetBulkSecret(ctx context.Context, storeName string, meta map[string]string) (data map[string]map[string]string, err error)

	// SaveState saves the raw data into store using default state options.
	SaveState(ctx context.Context, storeName, key string, data []byte, meta map[string]string, so ...StateOption) error

	// SaveBulkState saves multiple state item to store with specified options.
	SaveBulkState(ctx context.Context, storeName string, items ...*SetStateItem) error

	// GetState retrieves state from specific store using default consistency option.
	GetState(ctx context.Context, storeName, key string, meta map[string]string) (item *StateItem, err error)

	// GetStateWithConsistency retrieves state from specific store using provided state consistency.
	GetStateWithConsistency(ctx context.Context, storeName, key string, meta map[string]string, sc StateConsistency) (item *StateItem, err error)

	// GetBulkState retrieves state for multiple keys from specific store.
	GetBulkState(ctx context.Context, storeName string, keys []string, meta map[string]string, parallelism int32) ([]*BulkStateItem, error)

	// QueryStateAlpha1 runs a query against state store.
	QueryStateAlpha1(ctx context.Context, storeName, query string, meta map[string]string) (*QueryResponse, error)

	// DeleteState deletes content from store using default state options.
	DeleteState(ctx context.Context, storeName, key string, meta map[string]string) error

	// DeleteStateWithETag deletes content from store using provided state options and etag.
	DeleteStateWithETag(ctx context.Context, storeName, key string, etag *ETag, meta map[string]string, opts *StateOptions) error

	// ExecuteStateTransaction provides way to execute multiple operations on a specified store.
	ExecuteStateTransaction(ctx context.Context, storeName string, meta map[string]string, ops []*StateOperation) error

	// DeleteBulkState deletes content for multiple keys from store.
	DeleteBulkState(ctx context.Context, storeName string, keys []string, meta map[string]string) error

	// DeleteBulkState deletes content for multiple keys from store.
	DeleteBulkStateItems(ctx context.Context, storeName string, items []*DeleteStateItem) error

	// Shutdown the sidecar.
	Shutdown(ctx context.Context) error

	// WithTraceID adds existing trace ID to the outgoing context.
	WithTraceID(ctx context.Context, id string) context.Context

	// WithAuthToken sets Bhojpur Application API token on the instantiated client.
	WithAuthToken(token string)

	// Close cleans up all resources created by the Bhojpur Application runtime client.
	Close()

	// RegisterActorTimer registers an actor timer.
	RegisterActorTimer(ctx context.Context, req *RegisterActorTimerRequest) error

	// UnregisterActorTimer unregisters an actor timer.
	UnregisterActorTimer(ctx context.Context, req *UnregisterActorTimerRequest) error

	// RegisterActorReminder registers an actor reminder.
	RegisterActorReminder(ctx context.Context, req *RegisterActorReminderRequest) error

	// UnregisterActorReminder unregisters an actor reminder.
	UnregisterActorReminder(ctx context.Context, req *UnregisterActorReminderRequest) error

	// RenameActorReminder rename an actor reminder.
	RenameActorReminder(ctx context.Context, req *RenameActorReminderRequest) error

	// InvokeActor calls a method on an actor.
	InvokeActor(ctx context.Context, req *InvokeActorRequest) (*InvokeActorResponse, error)

	// GetActorState get actor state
	GetActorState(ctx context.Context, req *GetActorStateRequest) (data *GetActorStateResponse, err error)

	// SaveStateTransactionally save actor state
	SaveStateTransactionally(ctx context.Context, actorType, actorID string, operations []*ActorStateOperation) error

	// ImplActorClientStub is to impl user defined actor client stub
	ImplActorClientStub(actorClientStub actor.Client, opt ...config.Option)
}

// NewClient instantiates Bhojpur Application runtime client using APP_GRPC_PORT environment variable as port.
// Note, this default factory function creates Bhojpur Application client only once. All subsequent invocations
// will return the already created instance. To create multiple instances of the Bhojpur Application client,
// use one of the parameterized factory functions:
//   NewClientWithPort(port string) (client Client, err error)
//   NewClientWithAddress(address string) (client Client, err error)
//   NewClientWithConnection(conn *grpc.ClientConn) Client
//   NewClientWithSocket(socket string) (client Client, err error)
func NewClient() (client Client, err error) {
	port := os.Getenv(appPortEnvVarName)
	if port == "" {
		port = appPortDefault
	}
	var onceErr error
	doOnce.Do(func() {
		c, err := NewClientWithPort(port)
		onceErr = errors.Wrap(err, "error creating default Bhojpur Application client")
		defaultClient = c
	})

	return defaultClient, onceErr
}

// NewClientWithPort instantiates Bhojpur Application client using specific port.
func NewClientWithPort(port string) (client Client, err error) {
	if port == "" {
		return nil, errors.New("nil port")
	}
	return NewClientWithAddress(net.JoinHostPort("127.0.0.1", port))
}

// NewClientWithAddress instantiates Bhojpur Application client using specific address (including port).
func NewClientWithAddress(address string) (client Client, err error) {
	if address == "" {
		return nil, errors.New("nil address")
	}
	logger.Printf("Bhojpur Application client initializing for: %s", address)

	ctx, ctxCancel := context.WithTimeout(context.Background(), 1*time.Second)
	conn, err := grpc.DialContext(
		ctx,
		address,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		ctxCancel()
		return nil, errors.Wrapf(err, "error creating connection to '%s': %v", address, err)
	}
	if hasToken := os.Getenv(apiTokenEnvVarName); hasToken != "" {
		logger.Println("Bhojpur Application client uses API token")
	}

	return newClientWithConnectionAndCancelFunc(conn, ctxCancel), nil
}

// NewClientWithSocket instantiates Bhojpur Application client using specific socket.
func NewClientWithSocket(socket string) (client Client, err error) {
	if socket == "" {
		return nil, errors.New("nil socket")
	}
	logger.Printf("Bhojpur Application client initializing for: %s", socket)
	addr := fmt.Sprintf("unix://%s", socket)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrapf(err, "error creating connection to '%s': %v", addr, err)
	}
	if hasToken := os.Getenv(apiTokenEnvVarName); hasToken != "" {
		logger.Println("Bhojpur Application client uses API token")
	}
	return NewClientWithConnection(conn), nil
}

// NewClientWithConnection instantiates Bhojpur Application client using specific connection.
func NewClientWithConnection(conn *grpc.ClientConn) Client {
	return newClientWithConnectionAndCancelFunc(conn, func() {})
}

func newClientWithConnectionAndCancelFunc(
	conn *grpc.ClientConn,
	cancelFunc context.CancelFunc,
) Client {
	return &GRPCClient{
		connection:    conn,
		ctxCancelFunc: cancelFunc,
		protoClient:   pb.NewApplicationClient(conn),
		authToken:     os.Getenv(apiTokenEnvVarName),
	}
}

// GRPCClient is the gRPC implementation of Bhojpur Application client.
type GRPCClient struct {
	connection    *grpc.ClientConn
	ctxCancelFunc context.CancelFunc
	protoClient   pb.ApplicationClient
	authToken     string
	mux           sync.Mutex
}

// Close cleans up all resources created by the client.
func (c *GRPCClient) Close() {
	c.ctxCancelFunc()
	if c.connection != nil {
		c.connection.Close()
	}
}

// WithAuthToken sets Bhojpur Application API token on the instantiated client.
// Allows empty string to reset token on existing client.
func (c *GRPCClient) WithAuthToken(token string) {
	c.mux.Lock()
	c.authToken = token
	c.mux.Unlock()
}

// WithTraceID adds existing trace ID to the outgoing context.
func (c *GRPCClient) WithTraceID(ctx context.Context, id string) context.Context {
	if id == "" {
		return ctx
	}
	logger.Printf("using trace parent ID: %s", id)
	md := metadata.Pairs(traceparentKey, id)
	return metadata.NewOutgoingContext(ctx, md)
}

func (c *GRPCClient) withAuthToken(ctx context.Context) context.Context {
	if c.authToken == "" {
		return ctx
	}
	return metadata.NewOutgoingContext(ctx, metadata.Pairs(apiTokenKey, c.authToken))
}

// Shutdown the sidecar.
func (c *GRPCClient) Shutdown(ctx context.Context) error {
	_, err := c.protoClient.Shutdown(c.withAuthToken(ctx), &emptypb.Empty{})
	if err != nil {
		return errors.Wrap(err, "error shutting down the Bhojpur Application runtime sidecar")
	}
	return nil
}