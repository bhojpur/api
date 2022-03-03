package http

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
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bhojpur/api/pkg/service/common"
)

// AddBindingInvocationHandler appends provided binding invocation handler with its route to the service.
func (s *Server) AddBindingInvocationHandler(route string, fn common.BindingInvocationHandler) error {
	if route == "" {
		return fmt.Errorf("binding route required")
	}
	if fn == nil {
		return fmt.Errorf("binding handler required")
	}

	if !strings.HasPrefix(route, "/") {
		route = fmt.Sprintf("/%s", route)
	}

	s.mux.Handle(route, optionsHandler(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var content []byte
			if r.ContentLength > 0 {
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				content = body
			}

			// assuming Bhojpur Application runtime doesn't pass multiple values for key
			meta := map[string]string{}
			for k, values := range r.Header {
				// TODO: Need to figure out how to parse out only the headers set in the binding + Traceparent
				// if k == "raceparent" || strings.HasPrefix(k, "app") {
				for _, v := range values {
					meta[k] = v
				}
				// }
			}

			// execute handler
			in := &common.BindingEvent{
				Data:     content,
				Metadata: meta,
			}
			out, err := fn(r.Context(), in)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if out == nil {
				out = []byte("{}")
			}

			w.Header().Add("Content-Type", "application/json")
			if _, err := w.Write(out); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})))

	return nil
}
