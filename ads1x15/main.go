// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// ads1x15 reads from ADS1015/ADS1115 Analog-Digital Converters (ADC) via
// I²C interface.
package main

import (
	"fmt"
	"os"

	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/ads1x15"
	"periph.io/x/host/v3"
)

// Resistor values for voltage divider. ADC measures between r2 and ground.
// In this example a 24 tolerant cirquit is used with a voltage divider of
// r1=820kΩ and r2=120kΩ.
const (
	r1 = 820
	r2 = 120
)

func mainImpl() error {
	if _, err := host.Init(); err != nil {
		return err
	}
	bus, err := i2creg.Open("")
	if err != nil {
		return fmt.Errorf("failed to open I²C: %w", err)
	}
	defer bus.Close()
	adc, err := ads1x15.NewADS1015(bus, &ads1x15.DefaultOpts)
	if err != nil {
		return err
	}

	// Obtain an analog pin from the ADC.
	pin, err := adc.PinForChannel(ads1x15.Channel0, 1*physic.Volt, 1*physic.Hertz, ads1x15.SaveEnergy)
	if err != nil {
		return err
	}
	defer pin.Halt()

	// Read values from ADC.
	fmt.Println("Single reading")
	reading, err := pin.Read()
	if err != nil {
		return err
	}

	actualV := (reading.V * (r1 + r2) / r2)
	fmt.Println(actualV)

	// Read values continuously from ADC.
	fmt.Println("Continuous reading")
	c := pin.ReadContinuous()

	// TODO(maruel): Use os.Signal to gracefully shutdown on Ctrl-C.
	for reading := range c {
		actualV := (reading.V * (r1 + r2) / r2)
		fmt.Println(actualV)
	}
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "ads1x15: %s.\n", err)
		os.Exit(1)
	}
}
