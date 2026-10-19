// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// HTTP server shutdown compatibility shim for builds without Go 1.8 APIs.

//go:build !go1.8
// +build !go1.8

package main

func (s *webServer) Close() error {
	return nil
}
