# Bhojpur API - Distributed Computing Framework

The `core` library is a [Bhojpur.NET Platform](https://github.com/bhojpur/platform) application
programming interface for building high-performance, web-scale `applications` or `services`. It
is being applied by many core components to expose service capabilities within the framework.

| packages  | description                                                            |
|-----------|------------------------------------------------------------------------|
| common    | common protobuf definitions that are imported by multiple packages     |
| internals | gRPC and protobuf definitions, which is used for appication internal   |
| runtime   | application and callback services and its associated protobuf messages |
| operator  | application operator gRPC service                                      |
| placement | applicaion placement service                                           |
| sentry    | application sentry for CA service                                      |

## Client-side Protocol Buffer Generation

The `client-side` application programming interface is used to connect to core services provided by the
[Bhojpur.NET Platform](https://github.com/bhojpur/platform) ecosystem.

1. Install `protoc` version: [v3.14.0](https://github.com/protocolbuffers/protobuf/releases/tag/v3.14.0)

2. Install `protoc-gen-go` and `protoc-gen-go-grpc`

```bash
$ make init-proto
```

3. Generate gRPC protocol buffer Clients

```bash
$ make gen-proto
```

## Server-side Protocol Buffer Generation

The `server-side` application programming interfaces are used to implement to various core services
of the [Bhojpur.NET Platform](https://github.com/bhojpur/platform). The protocol buffer definitions
are open to allows alternatives. The [Bhojpur Application](https://github.com/bhojpur/application)
runtime implements the `server-side` to support a __cloud-native__ distributed computing platform.

## Update E2E Test Apps

Whenever there are breaking changes in the `.proto` files, we need to update the E2E test applications
to use the correct version of [Bhojpur Application](https://github.com/bhojpur/application) runtime
and its dependencies. It could be done by navigating to the tests folder and running the commands.

```bash
# Use the last commit of Bhojpur APIs and Applications
$ ./update_testapps_dependencies.sh be08e5520173beb93e5d5f047dbde405e78db658
```

**Note**: On Windows, use the `MingW` tools to execute the bash script

Check in all the `go.mod` files for the test applications that have now been modified to point
to the latest [Bhojpur Application](https://github.com/bhojpur/application) version.
