// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// hx711 reads from a 24-bits HX711 analog to digital converter.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/gpio/gpioutil"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/hx711"
	"periph.io/x/host/v3"
)

const timeout = time.Second

func mainImpl() error {
	clkPin := flag.String("clk", "", "Clock pin")
	dataPin := flag.String("data", "", "Data pin")
	gain := flag.Int("gain", 128,
		"Voltage gain. Must be one of 128, 64 or 32. Using 32 selects Channel B")
	cont := flag.Bool("cont", false, "Reads continuously from the ADC")
	usePollEdge := flag.Bool("poll-edge", false,
		"Poll the data pin instead of using edge detection")
	flag.Parse()

	if _, err := host.Init(); err != nil {
		return err
	}

	if *clkPin == "" {
		return fmt.Errorf("-clk is required")
	}
	if *dataPin == "" {
		return fmt.Errorf("-data is required")
	}

	clkPinReg := gpioreg.ByName(*clkPin)
	if clkPinReg == nil {
		return fmt.Errorf("clock pin %s can not be found", *clkPin)
	}
	dataPinReg := gpioreg.ByName(*dataPin)
	if dataPinReg == nil {
		return fmt.Errorf("data pin %s can not be found", *dataPin)
	}

	if *usePollEdge {
		dataPinReg = gpioutil.PollEdge(dataPinReg, 20*physic.KiloHertz)
	}

	dev, err := hx711.New(clkPinReg, dataPinReg)
	if err != nil {
		return err
	}

	switch *gain {
	case 128:
		if err := dev.SetInputMode(hx711.CHANNEL_A_GAIN_128); err != nil {
			return err
		}
	case 64:
		if err := dev.SetInputMode(hx711.CHANNEL_A_GAIN_64); err != nil {
			return err
		}
	case 32:
		if err := dev.SetInputMode(hx711.CHANNEL_B_GAIN_32); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid gain '%d', must be either 128, 64 or 32", *gain)
	}

	if *cont {
		ch := dev.ReadContinuous()
		for {
			fmt.Println(<-ch)
		}
	} else {
		value, err := dev.ReadTimeout(timeout)
		if err != nil {
			return err
		}
		fmt.Println(value)
	}
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "hx711: %s.\n", err)
		os.Exit(1)
	}
}
