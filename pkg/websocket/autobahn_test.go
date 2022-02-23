//go:build !js
// +build !js

package websocket_test

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
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bhojpur/api/pkg/websocket"
	"github.com/bhojpur/api/pkg/websocket/internal/errd"
	"github.com/bhojpur/api/pkg/websocket/internal/test/assert"
	"github.com/bhojpur/api/pkg/websocket/internal/test/wstest"
)

var excludedAutobahnCases = []string{
	// We skip the UTF-8 handling tests as there isn't any reason to reject invalid UTF-8, just
	// more performance overhead.
	"6.*", "7.5.1",

	// We skip the tests related to requestMaxWindowBits as that is unimplemented due
	// to limitations in compress/flate. See https://github.com/golang/go/issues/3155
	// Same with klauspost/compress which doesn't allow adjusting the sliding window size.
	"13.3.*", "13.4.*", "13.5.*", "13.6.*",
}

var autobahnCases = []string{"*"}

func TestAutobahn(t *testing.T) {
	t.Parallel()

	if os.Getenv("AUTOBAHN_TEST") == "" {
		t.SkipNow()
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*15)
	defer cancel()

	wstestURL, closeFn, err := wstestClientServer(ctx)
	assert.Success(t, err)
	defer closeFn()

	err = waitWS(ctx, wstestURL)
	assert.Success(t, err)

	cases, err := wstestCaseCount(ctx, wstestURL)
	assert.Success(t, err)

	t.Run("cases", func(t *testing.T) {
		for i := 1; i <= cases; i++ {
			i := i
			t.Run("", func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
				defer cancel()

				c, _, err := websocket.Dial(ctx, fmt.Sprintf(wstestURL+"/runCase?case=%v&agent=main", i), nil)
				assert.Success(t, err)
				err = wstest.EchoLoop(ctx, c)
				t.Logf("echoLoop: %v", err)
			})
		}
	})

	c, _, err := websocket.Dial(ctx, fmt.Sprintf(wstestURL+"/updateReports?agent=main"), nil)
	assert.Success(t, err)
	c.Close(websocket.StatusNormalClosure, "")

	checkWSTestIndex(t, "./ci/out/wstestClientReports/index.json")
}

func waitWS(ctx context.Context, url string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	for ctx.Err() == nil {
		c, _, err := websocket.Dial(ctx, url, nil)
		if err != nil {
			continue
		}
		c.Close(websocket.StatusNormalClosure, "")
		return nil
	}

	return ctx.Err()
}

func wstestClientServer(ctx context.Context) (url string, closeFn func(), err error) {
	serverAddr, err := unusedListenAddr()
	if err != nil {
		return "", nil, err
	}

	url = "ws://" + serverAddr

	specFile, err := tempJSONFile(map[string]interface{}{
		"url":           url,
		"outdir":        "ci/out/wstestClientReports",
		"cases":         autobahnCases,
		"exclude-cases": excludedAutobahnCases,
	})
	if err != nil {
		return "", nil, fmt.Errorf("failed to write spec: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*15)
	defer func() {
		if err != nil {
			cancel()
		}
	}()

	args := []string{"--mode", "fuzzingserver", "--spec", specFile,
		// Disables some server that runs as part of fuzzingserver mode.
		// See https://github.com/crossbario/autobahn-testsuite/blob/058db3a36b7c3a1edf68c282307c6b899ca4857f/autobahntestsuite/autobahntestsuite/wstest.py#L124
		"--webport=0",
	}
	wstest := exec.CommandContext(ctx, "wstest", args...)
	err = wstest.Start()
	if err != nil {
		return "", nil, fmt.Errorf("failed to start wstest: %w", err)
	}

	return url, func() {
		wstest.Process.Kill()
	}, nil
}

func wstestCaseCount(ctx context.Context, url string) (cases int, err error) {
	defer errd.Wrap(&err, "failed to get case count")

	c, _, err := websocket.Dial(ctx, url+"/getCaseCount", nil)
	if err != nil {
		return 0, err
	}
	defer c.Close(websocket.StatusInternalError, "")

	_, r, err := c.Reader(ctx)
	if err != nil {
		return 0, err
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, err
	}
	cases, err = strconv.Atoi(string(b))
	if err != nil {
		return 0, err
	}

	c.Close(websocket.StatusNormalClosure, "")

	return cases, nil
}

func checkWSTestIndex(t *testing.T, path string) {
	wstestOut, err := ioutil.ReadFile(path)
	assert.Success(t, err)

	var indexJSON map[string]map[string]struct {
		Behavior      string `json:"behavior"`
		BehaviorClose string `json:"behaviorClose"`
	}
	err = json.Unmarshal(wstestOut, &indexJSON)
	assert.Success(t, err)

	for _, tests := range indexJSON {
		for test, result := range tests {
			t.Run(test, func(t *testing.T) {
				switch result.BehaviorClose {
				case "OK", "INFORMATIONAL":
				default:
					t.Errorf("bad close behaviour")
				}

				switch result.Behavior {
				case "OK", "NON-STRICT", "INFORMATIONAL":
				default:
					t.Errorf("failed")
				}
			})
		}
	}

	if t.Failed() {
		htmlPath := strings.Replace(path, ".json", ".html", 1)
		t.Errorf("detected autobahn violation, see %q", htmlPath)
	}
}

func unusedListenAddr() (_ string, err error) {
	defer errd.Wrap(&err, "failed to get unused listen address")
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return "", err
	}
	l.Close()
	return l.Addr().String(), nil
}

func tempJSONFile(v interface{}) (string, error) {
	f, err := ioutil.TempFile("", "temp.json")
	if err != nil {
		return "", fmt.Errorf("temp file: %w", err)
	}
	defer f.Close()

	e := json.NewEncoder(f)
	e.SetIndent("", "\t")
	err = e.Encode(v)
	if err != nil {
		return "", fmt.Errorf("json encode: %w", err)
	}

	err = f.Close()
	if err != nil {
		return "", fmt.Errorf("close temp file: %w", err)
	}

	return f.Name(), nil
}
