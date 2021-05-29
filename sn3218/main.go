// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// TODO: This should be a periph-smoketest.

// sn3218 runs a smoketest on a SN3218 LED driver with 18 LEDs over an i2c bus.
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/devices/v3/sn3218"
	"periph.io/x/host/v3"
)

func mainImpl() error {
	if _, err := host.Init(); err != nil {
		return err
	}

	b, err := i2creg.Open("")
	if err != nil {
		return err
	}
	defer b.Close()

	d, err := sn3218.New(b)
	if err != nil {
		return err
	}
	defer d.Halt()

	if err = d.WakeUp(); err != nil {
		return fmt.Errorf("error enabling device: %w", err)
	}

	if err = d.BrightnessAll(1); err != nil {
		return fmt.Errorf("error setting brightness: %w", err)
	}

	// Switch LED 7 on
	if err = d.Switch(7, true); err != nil {
		return fmt.Errorf("error switching LED: %w", err)
	}
	time.Sleep(1000 * time.Millisecond)

	// Increase brightness for LED 7 to max
	if err = d.Brightness(7, 255); err != nil {
		return fmt.Errorf("error changing LED brightness: %w", err)
	}
	time.Sleep(1000 * time.Millisecond)

	// Get state of LED 7
	state, brightness, err := d.GetState(7)
	if err != nil {
		return fmt.Errorf("error reading LED state: %w", err)
	}
	log.Println("State: ", state, " - Brightness: ", brightness)

	// Switch all LEDs on
	if err = d.SwitchAll(true); err != nil {
		return fmt.Errorf("error switching all LEDs: %w", err)
	}
	time.Sleep(1000 * time.Millisecond)

	// Increase brightness for all
	if err = d.BrightnessAll(125); err != nil {
		return fmt.Errorf("error changing globalBrightness: %w", err)
	}
	time.Sleep(1000 * time.Millisecond)

	// Sleep mode to save energy, but keep state
	if err = d.Sleep(); err != nil {
		return fmt.Errorf("error disabling device: %w", err)
	}
	time.Sleep(1000 * time.Millisecond)

	// WakeUp again
	if err = d.WakeUp(); err != nil {
		return fmt.Errorf("error enabling device: %w", err)
	}
	time.Sleep(1000 * time.Millisecond)
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "sn3218: %s.\n", err)
		os.Exit(1)
	}
}
