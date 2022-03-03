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

// AddServiceInvocationHandler appends provided service invocation handler with its route to the service.
func (s *Server) AddServiceInvocationHandler(route string, fn common.ServiceInvocationHandler) error {
	if route == "" {
		return fmt.Errorf("service route required")
	}
	if fn == nil {
		return fmt.Errorf("invocation handler required")
	}

	if !strings.HasPrefix(route, "/") {
		route = fmt.Sprintf("/%s", route)
	}

	s.mux.Handle(route, optionsHandler(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if s.authToken != "" {
				token := r.Header.Get(common.APITokenKey)
				if token == "" || token != s.authToken {
					http.Error(w, "authentication failed.", http.StatusNonAuthoritativeInfo)
					return
				}
			}
			// capture http args
			e := &common.InvocationEvent{
				Verb:        r.Method,
				QueryString: r.URL.RawQuery,
				ContentType: r.Header.Get("Content-type"),
			}

			// check for post with no data
			if r.ContentLength > 0 {
				content, err := ioutil.ReadAll(r.Body)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				e.Data = content
			}

			// execute handler
			o, err := fn(r.Context(), e)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// write to response if handler returned data
			if o != nil && o.Data != nil {
				if o.ContentType != "" {
					w.Header().Set("Content-type", o.ContentType)
				}
				if _, err := w.Write(o.Data); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		})))

	return nil
}
