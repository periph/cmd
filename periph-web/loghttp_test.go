// Copyright 2026 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Tests logging wrapper support for HTTP response capabilities.

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoggingHandlerWithoutHijacker(t *testing.T) {
	t.Parallel()
	h := loggingHandler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if _, ok := w.(http.Hijacker); ok {
			t.Error("wrapper advertises hijacking when the underlying writer cannot hijack")
		}
	}))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil))
}
