// Copyright 2024 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.
//
// Functional smoke test for the host/gpioioctl package.
package gpiosmoketest

import (
	"errors"
	"flag"
	"log"
)

type SmokeTestT struct {
	error bool
	fatal bool
}

func (st *SmokeTestT) Error(args ...any) {
	st.error = true
	log.Println(args...)
}
func (st *SmokeTestT) Errorf(format string, args ...any) {
	st.error = true
	log.Printf(format, args...)
}
func (st *SmokeTestT) Fatal(args ...any) {
	st.fatal = true
	log.Fatal(args...)
}
func (st *SmokeTestT) Fatalf(format string, args ...any) {
	st.fatal = true
	log.Fatalf(format, args...)
}
func (st *SmokeTestT) Log(args ...any) {
	log.Println(args...)
}
func (st *SmokeTestT) Logf(format string, args ...any) {
	log.Printf(format, args...)
}
func (st *SmokeTestT) ErrorsOrFatals() bool {
	return st.error || st.fatal
}

type SmokeTest struct {
}

// Name implements periph-smoketest.SmokeTest.
func (s *SmokeTest) Name() string {
	return "gpio"
}

// Description implements periph-smoketest.SmokeTest.
func (s *SmokeTest) Description() string {
	return "Tests basic functionality, edge detection and input pull resistors"
}

// Run implements periph-smoketest.SmokeTest.
func (s *SmokeTest) Run(f *flag.FlagSet, args []string) error {
	st := &SmokeTestT{}
	TestWriteReadSinglePin(st)
	TestWaitForEdgeTimeout(st)
	TestWaitForEdgeSinglePin(st)
	TestHalt(st)
	TestLineSetCreation(st)
	TestLineSetReadWrite(st)
	TestLineSetWaitForEdgeTimeout(st)
	TestLineSetHalt(st)
	TestLineSetWaitForEdge(st)
	TestLineSetConfigWithOverride(st)
	if st.ErrorsOrFatals() {
		return errors.New("Smoketest failure.")
	}
	return nil
}
