package gpiosmoketest

// Copyright 2024 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.
//
// The In/Out tests depend upon having a jumper wire connecting _OUT_LINE and
// _IN_LINE

import (
	"time"

	"periph.io/x/conn/v3/driver/driverreg"
	"periph.io/x/conn/v3/gpio"
    "periph.io/x/host/v3/gpioioctl"
)

const (
	_OUT_LINE = "GPIO5"
	_IN_LINE  = "GPIO13"
)

func init() {
	_, _ = driverreg.Init()
}

func TestWriteReadSinglePin(t *SmokeTestT) {
	var err error
	chip := gpioioctl.Chips[0]
	inLine := chip.ByName(_IN_LINE)
	outLine := chip.ByName(_OUT_LINE)
	defer inLine.Close()
	defer outLine.Close()
	err = outLine.Out(true)
	if err != nil {
		t.Errorf("outLine.Out() %s", err)
	}
	if val := inLine.Read(); !val {
		t.Error("Error reading/writing GPIO Pin. Expected true, received false!")
	}
	if inLine.Pull() != gpio.PullUp {
		t.Errorf("Pull() returned %s expected %s", gpioioctl.PullLabels[inLine.Pull()], gpioioctl.PullLabels[gpio.PullUp])
	}
	err = outLine.Out(false)
	if err != nil {
		t.Errorf("outLine.Out() %s", err)
	}
	if val := inLine.Read(); val {
		t.Error("Error reading/writing GPIO Pin. Expected false, received true!")
	}
	/*
		By Design, lines should auto change directions if Read()/Out() are called
		and they don't match.
	*/
	err = inLine.Out(false)
	if err != nil {
		t.Errorf("inLine.Out() %s", err)
	}
	time.Sleep(500 * time.Millisecond)
	err = inLine.Out(true)
	if err != nil {
		t.Errorf("inLine.Out() %s", err)
	}
	if val := outLine.Read(); !val {
		t.Error("Error read/writing with auto-reverse of line functions.")
	}
	err = inLine.Out(false)
	if err != nil {
		t.Errorf("TestWriteReadSinglePin() %s", err)
	}
	if val := outLine.Read(); val {
		t.Error("Error read/writing with auto-reverse of line functions.")
	}

}

func clearEdges(line gpio.PinIn) bool {
	result := false
	for line.WaitForEdge(10 * time.Millisecond) {
		result = true
	}
	return result
}

func TestWaitForEdgeTimeout(t *SmokeTestT) {
	line := gpioioctl.Chips[0].ByName(_IN_LINE)
	defer line.Close()
	err := line.In(gpio.PullUp, gpio.BothEdges)
	if err != nil {
		t.Error(err)
	}
	clearEdges(line)
	tStart := time.Now().UnixMilli()
	line.WaitForEdge(5 * time.Second)
	tEnd := time.Now().UnixMilli()
	tDiff := tEnd - tStart
	if tDiff < 4500 || tDiff > 5500 {
		t.Errorf("timeout duration failure. Expected duration: 5000, Actual duration: %d", tDiff)
	}
}

// Test detection of rising, falling, and both.
func TestWaitForEdgeSinglePin(t *SmokeTestT) {
	tests := []struct {
		startVal gpio.Level
		edge     gpio.Edge
		writeVal gpio.Level
	}{
		{startVal: false, edge: gpio.RisingEdge, writeVal: true},
		{startVal: true, edge: gpio.FallingEdge, writeVal: false},
		{startVal: false, edge: gpio.BothEdges, writeVal: true},
		{startVal: true, edge: gpio.BothEdges, writeVal: false},
	}
	var err error
	line := gpioioctl.Chips[0].ByName(_IN_LINE)
	outLine := gpioioctl.Chips[0].ByName(_OUT_LINE)
	defer line.Close()
	defer outLine.Close()

	for _, test := range tests {
		err = outLine.Out(test.startVal)
		if err != nil {
			t.Errorf("set initial value. %s", err)
		}
		err = line.In(gpio.PullUp, test.edge)
		if err != nil {
			t.Errorf("line.In() %s", err)
		}
		clearEdges(line)
		err = outLine.Out(test.writeVal)
		if err != nil {
			t.Errorf("outLine.Out() %s", err)
		}
		if edgeReceived := line.WaitForEdge(time.Second); !edgeReceived {
			t.Errorf("Expected Edge %s was not received on transition from %t to %t", gpioioctl.EdgeLabels[test.edge], test.startVal, test.writeVal)
		}
	}
}

func TestHalt(t *SmokeTestT) {
	line := gpioioctl.Chips[0].ByName(_IN_LINE)
	defer line.Close()
	err := line.In(gpio.PullUp, gpio.BothEdges)
	if err != nil {
		t.Fatalf("TestHalt() %s", err)
	}
	clearEdges(line)
	// So what we'll do here is setup a goroutine to wait three seconds and then send a halt.
	go func() {
		time.Sleep(time.Second * 3)
		err = line.Halt()
		if err != nil {
			t.Error(err)
		}
	}()
	tStart := time.Now().UnixMilli()
	line.WaitForEdge(time.Second * 30)
	tEnd := time.Now().UnixMilli()
	tDiff := tEnd - tStart
	if tDiff > 3500 {
		t.Errorf("error calling halt to interrupt WaitForEdge() Duration %d exceeded expected value.", tDiff)
	}
}
