// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Read web assets from disk for debug builds.

//go:build debug
// +build debug

package main

import (
	"log"
	"os"
)

const cacheControlContent = "Cache-Control:no-cache,private"

// getContent reads an asset path selected by the handlers, never a request URL.
func getContent(path string) []byte {
	r, err := os.ReadFile(path)
	if err != nil {
		log.Printf("getContent(%q): %v", path, err)
	}
	return r
}
