// Copyright 2026 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Tests web server startup, authentication, and JSON request validation.

package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestWebServer(t *testing.T) {
	t.Run("cookie", func(t *testing.T) {
		s := &webServer{}
		w := httptest.NewRecorder()
		token := s.setXSRFCookie("127.0.0.1:1234", w)
		cookies := w.Result().Cookies()
		if len(cookies) != 1 {
			t.Fatalf("cookies = %v", cookies)
		}
		c := cookies[0]
		if !c.HttpOnly || c.SameSite != http.SameSiteStrictMode || c.Path != "/" {
			t.Fatalf("cookie protection = %s", c.String())
		}
		if c.Secure {
			t.Fatal("plain HTTP token cookie must remain usable")
		}
		if !s.validateToken(token, "127.0.0.1") {
			t.Fatal("fresh token was rejected")
		}
	})
	t.Run("malformed JSON", func(t *testing.T) {
		for _, body := range []string{`{"value":`, `{"unexpected":true}`, `{"unexpected":"data"}`, `{}` + `{}`} {
			t.Run(body, func(t *testing.T) {
				s := &webServer{}
				called := false
				h := s.api(func(in struct{ Value string }) (string, int) { called = true; return "ok", http.StatusOK })
				r := httptest.NewRequest(http.MethodPost, "/api", strings.NewReader(body))
				r.RemoteAddr = "127.0.0.1:1234"
				r.Header.Set("Content-Type", "application/json")
				r.AddCookie(&http.Cookie{Name: "XSRF-TOKEN", Value: s.generateToken("127.0.0.1", time.Now())})
				w := httptest.NewRecorder()
				h(w, r)
				if w.Code != http.StatusBadRequest || called {
					t.Fatalf("status=%d handlerCalled=%v", w.Code, called)
				}
			})
		}
	})
}

func TestAPIValidJSON(t *testing.T) {
	for _, withArg := range []bool{false, true} {
		s := &webServer{}
		called := false
		var handler any = func() (string, int) { called = true; return "ok", http.StatusOK }
		if withArg {
			handler = func(in struct{ Value string }) (string, int) { called = true; return in.Value, http.StatusOK }
		}
		body := `{}`
		if withArg {
			body = `{"Value":"ok"}`
		}
		r := httptest.NewRequest(http.MethodPost, "/api", strings.NewReader(body+" \n"))
		r.RemoteAddr = "127.0.0.1:1234"
		r.Header.Set("Content-Type", "application/json")
		r.AddCookie(&http.Cookie{Name: "XSRF-TOKEN", Value: s.generateToken("127.0.0.1", time.Now())})
		w := httptest.NewRecorder()
		s.api(handler)(w, r)
		if w.Code != http.StatusOK || !called {
			t.Fatalf("status=%d handlerCalled=%v", w.Code, called)
		}
	}
}

func TestAPINoArgsRejectsData(t *testing.T) {
	s := &webServer{}
	called := false
	h := s.api(func() (string, int) { called = true; return "ok", http.StatusOK })
	r := httptest.NewRequest(http.MethodPost, "/api", strings.NewReader(`{"unexpected":"data"}`))
	r.RemoteAddr = "127.0.0.1:1234"
	r.Header.Set("Content-Type", "application/json")
	r.AddCookie(&http.Cookie{Name: "XSRF-TOKEN", Value: s.generateToken("127.0.0.1", time.Now())})
	w := httptest.NewRecorder()
	h(w, r)
	if w.Code != http.StatusBadRequest || called {
		t.Fatalf("status=%d handlerCalled=%v", w.Code, called)
	}
}
