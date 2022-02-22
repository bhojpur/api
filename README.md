# Bhojpur API - Platform Access Library

The Bhojpur API is a standard `client-side` library for accessing the `server-side` of the [Bhojpur.NET Platform](https://github.com/bhojpur/platform) ecosystem. It offers a comprehensive collection of standards compliant, web-based `application programming interface` to be able to utilize a wide range of services (i.e., available in a __self-hosted__ or __managed hosting__ model). Also, it is available in several __programming languages__ (e.g., Go, Javascript, Python).

Some of our foundation frameworks (e.g. [Bhojpur GUI](https://github.com/bhojpur/gui) or [Bhojpur Application](https://github.com/bhojpur/application)) levevage these [APIs](https://github.com/bhojpur/api) for building web-scale `applications` and/or `services`.

Being a `client-side` library, it compiles without any compilation dependency on the services it tries to access. It is assumed that either HTTP or HTTP/S access should be sufficient.

## Bhojpur API - WebGL Access

Firstly, you should install `gopherjs` and verify the system using the following commands.

```sh
$ go get -u github.com/gopherjs/gopherjs
$ gopherjs version
```

Then, try compiling the sample program by issuing following `gopherjs` command.

```sh
$ gopherjs build internal/web/main.go
```

You need a basic `web server` to run serve the sample files (e.g., .html, .js) using
a standard `web browsers` (e.g., Chrome, Firefox). Therefore, issue the following commands.

```sh
$ npm install --global http-server
$ cd internal/web
$ http-server
```

Now, open a `web browser` tab and point address to `http://localhost:8080` to see the results.
