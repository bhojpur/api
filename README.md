# Bhojpur API - Platform Access Library

The Bhojpur API is a standard `client-side` library for accessing the `server-side` of
the [Bhojpur.NET Platform](https://github.com/bhojpur/platform) ecosystem. It offers a
comprehensive collection of standards compliant, web-based `application programming interface`
to be able to utilize a wide range of services (i.e., available in a __self-hosted__ or
__managed hosting__ model). Also, it is available in several __programming languages__
(e.g., Go, Javascript, Python).

Some of our foundation frameworks (e.g. [Bhojpur GUI](https://github.com/bhojpur/gui) or
[Bhojpur Application](https://github.com/bhojpur/application)) levevage these
[APIs](https://github.com/bhojpur/api) for building web-scale `applications` and/or `services`.

Being a `client-side` library, it compiles without any compilation dependency on the
services it tries to access. It is assumed that either HTTP or HTTP/S access should be
sufficient.

## Bhojpur API - Protocol Buffer Generation

The application programming interface is defined using protocol buffers. You could generate
the client stubs / server skeletons using the following commands.

```bash
$ make init-proto
$ make gen-proto
```

Please refer to [core](https://github.com/bhojpur/api/tree/main/pkg/core) library for APIs.

## Bhojpur API - WebSocket Access Library

The `websocket` client/server library is available in __./pkg/websocket__ folder. It supports
assess using `Javascript` and/or `WebAssembly` for the web.

## Bhojpur API - WebGL Access Library

Firstly, you should install `gopherjs` and verify the system using the following commands.

```bash
$ go get -u github.com/gopherjs/gopherjs
$ gopherjs version
```

Then, try compiling a sample program by issuing the following `gopherjs` command.

```bash
$ gopherjs build internal/web/main.go
```

You need a basic `web server` to run serve the sample files (e.g., .html, .js) using
a standard `web browsers` (e.g., Chrome, Firefox). Therefore, issue the following commands.

```bash
$ npm install --global http-server
$ cd internal/web
$ http-server
```

Now, open a `web browser` tab and point address to `http://localhost:8080` to see the results.

## Bhojpur API - Client Library

The `Client` library helps you build [Bhojpur Application](https://github.com/bhojpur/application).
It supports all public [Bhojpur APIs](https://docs.bhojpur/reference/api/), while focusing on
idiomatic Go experience and developer productivity.

### Simple Usage

> Assuming you already have [installed](https://golang.org/doc/install) `Go`

The Bhojpur API client includes two packages: `client` (for invoking public Bhojpur Application
runtime APIs), and `service` (to create services that will be invoked by Bhojpur Application
runtime, this is sometimes referred to as `callback`).

### Creating a Client

import Bhojpur API `client` package:

```go
import "github.com/bhojpur/api/pkg/client"
```

#### Quick start

```go
package main

import (
    app "github.com/bhojpur/api/pkg/client"
)

func main() {
    client, err := app.NewClient()
    if err != nil {
        panic(err)
    }
    defer client.Close()
    // TODO: use the client here, see below for examples 
}
```

Assuming you have [Bhojpur Application CLI](https://docs.bhojpur.net/getting-started/install-app/)
installed, you can then launch your application locally like this:

```shell
$ appctl run --app-id example-service \
            --app-protocol grpc \
            --app-port 50001 \
            go run main.go
```

#### Usage

The `client` library supports all the building blocks exposed by Bhojpur Application runtime API.
Let's review these one by one:

##### State

For simple use-cases, Bhojpur API client provides easy to use `Save`, `Get`, and `Delete` methods:

```go
ctx := context.Background()
data := []byte("hello")
store := "my-store" // defined in the component YAML 

// save state with the key key1, default options: strong, last-write
if err := client.SaveState(ctx, store, "key1", data, nil); err != nil {
    panic(err)
}

// get state for key key1
item, err := client.GetState(ctx, store, "key1", nil)
if err != nil {
    panic(err)
}
fmt.Printf("data [key:%s etag:%s]: %s", item.Key, item.Etag, string(item.Value))

// delete state for key key1
if err := client.DeleteState(ctx, store, "key1", nil); err != nil {
    panic(err)
}
```

For more granular control, the Bhojpur API client exposes `SetStateItem` type, which can be used
to gain more control over the state operations and allow for multiple items to be saved at once:

```go
item1 := &app.SetStateItem{
    Key:  "key1",
    Etag: &ETag{
        Value: "1",
    },
    Metadata: map[string]string{
        "created-on": time.Now().UTC().String(),
    },
    Value: []byte("hello"),
    Options: &app.StateOptions{
        Concurrency: app.StateConcurrencyLastWrite,
        Consistency: app.StateConsistencyStrong,
    },
}

item2 := &app.SetStateItem{
    Key:  "key2",
    Metadata: map[string]string{
        "created-on": time.Now().UTC().String(),
    },
    Value: []byte("hello again"),
}

item3 := &app.SetStateItem{
    Key:  "key3",
    Etag: &app.ETag{
	Value: "1",
    },
    Value: []byte("hello again"),
}

if err := client.SaveBulkState(ctx, store, item1, item2, item3); err != nil {
    panic(err)
}
```

Similarly, `GetBulkState` method provides a way to retrieve multiple state items in a single operation:

```go
keys := []string{"key1", "key2", "key3"}
items, err := client.GetBulkState(ctx, store, keys, nil,100)
```

And, the `ExecuteStateTransaction` method to execute multiple `upsert` or `delete` operations
transactionally.

```go
ops := make([]*app.StateOperation, 0)

op1 := &app.StateOperation{
    Type: app.StateOperationTypeUpsert,
    Item: &app.SetStateItem{
        Key:   "key1",
        Value: []byte(data),
    },
}
op2 := &app.StateOperation{
    Type: app.StateOperationTypeDelete,
    Item: &app.SetStateItem{
        Key:   "key2",
    },
}
ops = append(ops, op1, op2)
meta := map[string]string{}
err := testClient.ExecuteStateTransaction(ctx, store, meta, ops)
```

##### PubSub

To publish data onto a topic, the Bhojpur API client provides a simple method:

```go
data := []byte(`{ "id": "a123", "value": "abcdefg", "valid": true }`)
if err := client.PublishEvent(ctx, "component-name", "topic-name", data); err != nil {
    panic(err)
}
```

##### Service Invocation

To invoke a specific method on another service running with Bhojpur Application runtime sidecar,
the Bhojpur API client provides two options. To invoke a service without any data:

```go
resp, err := client.InvokeMethod(ctx, "app-id", "method-name", "post")
``` 

And, to invoke a service with data:

```go
content := &app.DataContent{
    ContentType: "application/json",
    Data:        []byte(`{ "id": "a123", "value": "demo", "valid": true }`),
}

resp, err = client.InvokeMethodWithContent(ctx, "app-id", "method-name", "post", content)
```

##### Bindings

Similar to the Service, the Bhojpur API client provides two methods to invoke an operation on a
[Bhojpur Application defined binding](https://docs.bhojpur.net/developing-applications/building-blocks/bindings/). The Bhojpur Application runtime supports input, output, and bidirectional bindings.

For simple, output only binding:

```go
in := &app.InvokeBindingRequest{ Name: "binding-name", Operation: "operation-name" }
err = client.InvokeOutputBinding(ctx, in)
```

To invoke method with content and metadata:

```go
in := &app.InvokeBindingRequest{
    Name:      "binding-name",
    Operation: "operation-name",
    Data: []byte("hello"),
    Metadata: map[string]string{"k1": "v1", "k2": "v2"},
}

out, err := client.InvokeBinding(ctx, in)
```

##### Secrets

The Bhojpur API client also provides access to the runtime secrets that can be backed by any number of
secret stores (e.g. Kubernetes Secrets, HashiCorp Vault, or Azure KeyVault):

```go
opt := map[string]string{
    "version": "2",
}

secret, err := client.GetSecret(ctx, "store-name", "secret-name", opt)
```

##### Authentication

By default, the Bhojpur Application runtime relies on the network boundary to limit access to its
API. If however the target Bhojpur Application runtime API is configured with token-based authentication,
users can configure the Bhojpur API client with that token in two ways:

###### Environment Variable

If the `APP_API_TOKEN` environment variable is defined, the Bhojpur Application runtime will automatically
use it to augment its Bhojpur Application runtime API invocations to ensure authentication.

###### Explicit Method

In addition, users can also set the API token explicitly on any Bhojpur API client instance. This approach
is helpful in cases when the user code needs to create multiple clients for different Bhojpur Application
runtime API endpoints.

```go
func main() {
    client, err := app.NewClient()
    if err != nil {
        panic(err)
    }
    defer client.Close()
    client.WithAuthToken("your-app-API-token-here")
}
```

### Service (callback)

In addition to the `client` capabilities that allow you to call into the Bhojpur Application runtime API,
the Bhojpur API also provides `service` package to help you bootstrap Bhojpur Application runtime callback
services in either gRPC or HTTP. Instructions on how to use it are located [here](./pkg/service/README.md)
