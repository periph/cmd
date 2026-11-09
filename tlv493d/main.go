// Copyright 2020 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// tlv493d measures the magnetic flux sensed by a tlv493d component.
package main

import (
	"fmt"
	"log"

	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/tlv493d"
	"periph.io/x/host/v3"
)

func run() error {
	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		return err
	}

	// Open default I²C bus.
	bus, err := i2creg.Open("")
	if err != nil {
		return fmt.Errorf("failed to open I²C: %w", err)
	}
	defer func() {
		if closeErr := bus.Close(); closeErr != nil {
			log.Printf("Close bus: %v", closeErr)
		}
	}()

	// Create a new TLV493D hall effect sensor.
	tlv, err := tlv493d.New(bus, &tlv493d.DefaultOpts)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := tlv.Halt(); closeErr != nil {
			log.Printf("Halt tlv: %v", closeErr)
		}
	}()

	// Read a single value.
	if err = tlv.SetMode(tlv493d.LowPowerMode); err != nil {
		return err
	}
	fmt.Println("Single reading")
	reading, err := tlv.Read(tlv493d.HighPrecisionWithTemperature)

	if err != nil {
		return err
	}

	fmt.Println(reading)

	// Read values continuously from the sensor.
	fmt.Println("Continuous reading")
	c, err := tlv.ReadContinuous(100*physic.Hertz, tlv493d.LowPrecision)
	if err != nil {
		return err
	}

	for reading := range c {
		fmt.Println(reading)
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
