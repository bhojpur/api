# Bhojpur API - Basic Actor

## Steps

### Prepare

- [Bhojpur Application](https://github.com/bhojpur/application) `runtime` installed
- [Bhojpur Application](https://github.com/bhojpur/application) `placement` server started
- `redis` installed to function as a state store

### Run Actor Server

<!-- STEP
name: Run Actor server
output_match_mode: substring
expected_stdout_lines:
  - '== MANAGED APPLICATION == call get user req = &{abc 123}'
  - '== MANAGED APPLICATION == get req = pramila'
  - '== MANAGED APPLICATION == get post request = pramila'
  - '== MANAGED APPLICATION == get req = hello'
  - '== MANAGED APPLICATION == get req = hello'
  - '== MANAGED APPLICATION == receive reminder = testReminderName state = "hello" duetime = 5s period = 5s'
  - '== MANAGED APPLICATION == receive reminder = testReminderName state = "hello" duetime = 5s period = 5s'
background: true
sleep: 30
-->

```bash
$ cd internal/basic
$ appctl run --app-id basic-server \
            --app-protocol http \
            --app-port 3000 \
            --app-http-port 3500 \
            --log-level debug \
            --components-path ./config \
            go run ./server/main.go
```

<!-- END_STEP -->

### Run Actor Client

<!-- STEP
name: Run Actor Client
output_match_mode: substring
expected_stdout_lines:
  - '== MANAGED APPLICATION == get user result = &{abc 123}'
  - '== MANAGED APPLICATION == get invoke result = pramila'
  - '== MANAGED APPLICATION == get post result = pramila'
  - '== MANAGED APPLICATION == get result = get result'
  - '== MANAGED APPLICATION == start timer'
  - '== MANAGED APPLICATION == stop timer'
  - '== MANAGED APPLICATION == start reminder'
  - '== MANAGED APPLICATION == stop reminder'
  - '== MANAGED APPLICATION == get user = {Name: Age:1}'
  - '== MANAGED APPLICATION == get user = {Name: Age:2}'

background: true
sleep: 40
-->

```bash
$ appctl run --app-id basic-client \
            --log-level debug \
            --components-path ./config \
            go run ./client/main.go
```

<!-- END_STEP -->

### Cleanup

<!-- STEP
expected_stdout_lines: 
  - '✅  app stopped successfully: basic-server'
expected_stderr_lines:
name: Shutdown the Bhojpur Application runtime successfully
-->

```bash
$ appctl stop --app-id  basic-server
(lsof -i:3000 | grep main) | awk '{print $2}' | xargs  kill
```

<!-- END_STEP -->

## Result

- **client side**

```
== MANAGED APPLICATION == Bhojpur Application client initializing for: 127.0.0.1:55776
== MANAGED APPLICATION == get user result = &{abc 123}
== MANAGED APPLICATION == get invoke result = pramila
== MANAGED APPLICATION == get post result = pramila
== MANAGED APPLICATION == get result = get result
== MANAGED APPLICATION == start timer
== MANAGED APPLICATION == stop timer
== MANAGED APPLICATION == start reminder
== MANAGED APPLICATION == stop reminder
== MANAGED APPLICATION == get user = {Name: Age:1}
== MANAGED APPLICATION == get user = {Name: Age:2}
✅  Exited the Bhojpur Application runtime successfully

```

- **server-side**

```bash
== MANAGED APPLICATION == call get user req = &{abc 123}
== MANAGED APPLICATION == get req = pramila
== MANAGED APPLICATION == get post request = pramila
== MANAGED APPLICATION == get req = hello
== MANAGED APPLICATION == get req = hello
== MANAGED APPLICATION == receive reminder = testReminderName state = "hello" duetime = 5s period = 5s
== MANAGED APPLICATION == receive reminder = testReminderName state = "hello" duetime = 5s period = 5s
```
