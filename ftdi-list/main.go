// Copyright 2017 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// ftdi prints out information about the FTDI devices found on the USB bus.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"periph.io/x/d2xx"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/ftdi"
)

func process(d ftdi.Dev) {
	i := ftdi.Info{}
	d.Info(&i)
	fmt.Printf("  Type:           %s\n", i.Type)
	fmt.Printf("  Vendor ID:      %#04x\n", i.VenID)
	fmt.Printf("  Device ID:      %#04x\n", i.DevID)

	ee := ftdi.EEPROM{}
	if err := d.EEPROM(&ee); err == nil {
		fmt.Printf("  Manufacturer:   %s\n", ee.Manufacturer)
		fmt.Printf("  ManufacturerID: %s\n", ee.ManufacturerID)
		fmt.Printf("  Desc:           %s\n", ee.Desc)
		fmt.Printf("  Serial:         %s\n", ee.Serial)

		h := ee.AsHeader()
		fmt.Printf("  MaxPower:       %dmA\n", h.MaxPower)
		fmt.Printf("  SelfPowered:    %x\n", h.SelfPowered)
		fmt.Printf("  RemoteWakeup:   %x\n", h.RemoteWakeup)
		fmt.Printf("  PullDownEnable: %x\n", h.PullDownEnable)
		switch i.Type {
		case "FT232H":
			p := ee.AsFT232H()
			fmt.Printf("  CSlowSlew:      %t\n", p.ACSlowSlew != 0)        // AC bus pins have slow slew
			fmt.Printf("  CSchmittInput:  %t\n", p.ACSchmittInput != 0)    // AC bus pins are Schmitt input
			fmt.Printf("  CDriveCurrent:  %dmA\n", p.ACDriveCurrent*2)     // Valid values are 4mA, 8mA, 12mA, 16mA in 2mA units
			fmt.Printf("  DSlowSlew:      %t\n", p.ADSlowSlew != 0)        // AD bus pins have slow slew
			fmt.Printf("  DSchmittInput:  %t\n", p.ADSchmittInput != 0)    // AD bus pins are Schmitt input
			fmt.Printf("  DDriveCurrent:  %dmA\n", p.ADDriveCurrent)       // Valid values are 4mA, 8mA, 12mA, 16mA in 2mA units
			fmt.Printf("  Cbus0:          %s\n", p.Cbus0)                  //
			fmt.Printf("  Cbus1:          %s\n", p.Cbus1)                  //
			fmt.Printf("  Cbus2:          %s\n", p.Cbus2)                  //
			fmt.Printf("  Cbus3:          %s\n", p.Cbus3)                  //
			fmt.Printf("  Cbus4:          %s\n", p.Cbus4)                  //
			fmt.Printf("  Cbus5:          %s\n", p.Cbus5)                  //
			fmt.Printf("  Cbus6:          %s\n", p.Cbus6)                  //
			fmt.Printf("  Cbus7:          %s\n", p.Cbus7)                  // C7 is limited a sit can only do 'suspend on C7 low'. Defaults pull down.
			fmt.Printf("  Cbus8:          %s\n", p.Cbus8)                  //
			fmt.Printf("  Cbus9:          %s\n", p.Cbus9)                  //
			fmt.Printf("  FT1248Cpol:     %t\n", p.FT1248Cpol != 0)        // FT1248 clock polarity - clock idle high (true) or clock idle low (false)
			fmt.Printf("  FT1248Lsb:      %t\n", p.FT1248Lsb != 0)         // FT1248 data is LSB (true), or MSB (false)
			fmt.Printf("  FT1248FlowCtrl: %t\n", p.FT1248FlowControl != 0) // FT1248 flow control enable
			fmt.Printf("  IsFifo:         %t\n", p.IsFifo != 0)            // Interface is 245 FIFO
			fmt.Printf("  IsFifoTar:      %t\n", p.IsFifoTar != 0)         // Interface is 245 FIFO CPU target
			fmt.Printf("  IsFastSer:      %t\n", p.IsFastSer != 0)         // Interface is Fast serial
			fmt.Printf("  IsFT1248:       %t\n", p.IsFT1248 != 0)          // Interface is FT1248
			fmt.Printf("  PowerSaveEnabl: %t\n", p.PowerSaveEnable != 0)   // Suspect on ACBus7 low
			fmt.Printf("  DriverType:     %d\n", p.DriverType)             // 0 is D2XX, 1 is VCP
		case "FT2232H":
			p := ee.AsFT2232H()
			fmt.Printf("  ALSlowSlew:      %t\n", p.ALSlowSlew != 0)      // AL pins have slow slew
			fmt.Printf("  ALSchmittInput:  %t\n", p.ALSlowSlew != 0)      // AL pins are Schmitt input
			fmt.Printf("  ALDriveCurrent:  %dmA\n", p.ALDriveCurrent*2)   // Valid values are 4mA, 8mA, 12mA, 16mA in 2mA units
			fmt.Printf("  AHSlowSlew:      %t\n", p.AHSlowSlew != 0)      // AH pins have slow slew
			fmt.Printf("  AHSchmittInput:  %t\n", p.AHSlowSlew != 0)      // AH pins are Schmitt input
			fmt.Printf("  AHDriveCurrent:  %dmA\n", p.AHDriveCurrent*2)   // Valid values are 4mA, 8mA, 12mA, 16mA in 2mA units
			fmt.Printf("  BLSlowSlew:      %t\n", p.ALSlowSlew != 0)      // BL pins have slow slew
			fmt.Printf("  BLSchmittInput:  %t\n", p.ALSlowSlew != 0)      // BL pins are Schmitt input
			fmt.Printf("  BLDriveCurrent:  %dmA\n", p.ALDriveCurrent*2)   // Valid values are 4mA, 8mA, 12mA, 16mA in 2mA units
			fmt.Printf("  BHSlowSlew:      %t\n", p.AHSlowSlew != 0)      // BH pins have slow slew
			fmt.Printf("  BHSchmittInput:  %t\n", p.AHSlowSlew != 0)      // BH pins are Schmitt input
			fmt.Printf("  BHDriveCurrent:  %dmA\n", p.AHDriveCurrent*2)   // Valid values are 4mA, 8mA, 12mA, 16mA in 2mA units
			fmt.Printf("  AIsFifo:         %t\n", p.AIsFifo != 0)         // Interface is 245 FIFO
			fmt.Printf("  AIsFifoTar:      %t\n", p.AIsFifoTar != 0)      // Interface is 245 FIFO CPU target
			fmt.Printf("  AIsFastSer:      %t\n", p.AIsFastSer != 0)      // Interface is Fast serial
			fmt.Printf("  BIsFifo:         %t\n", p.BIsFifo != 0)         // Interface is 245 FIFO
			fmt.Printf("  BIsFifoTar:      %t\n", p.BIsFifoTar != 0)      // Interface is 245 FIFO CPU target
			fmt.Printf("  BIsFastSer:      %t\n", p.BIsFastSer != 0)      // Interface is Fast serial
			fmt.Printf("  PowerSaveEnable: %t\n", p.PowerSaveEnable != 0) // Using BCBUS7 to save power for self-powered designs
			fmt.Printf("  ADriverType:     %t\n", p.ADriverType != 0)     // 0 is D2XX, 1 is VCP
			fmt.Printf("  BDriverType:     %t\n", p.BDriverType != 0)     // 0 is D2XX, 1 is VCP
		case "FT232R":
			p := ee.AsFT232R()
			fmt.Printf("  IsHighCurrent:  %t\n", p.IsHighCurrent != 0) // High Drive I/Os; 3mA instead of 1mA (@3.3V)
			fmt.Printf("  UseExtOsc:      %t\n", p.UseExtOsc != 0)     // Use external oscillator
			fmt.Printf("  InvertTXD:      %t\n", p.InvertTXD != 0)     //
			fmt.Printf("  InvertRXD:      %t\n", p.InvertRXD != 0)     //
			fmt.Printf("  InvertRTS:      %t\n", p.InvertRTS != 0)     //
			fmt.Printf("  InvertCTS:      %t\n", p.InvertCTS != 0)     //
			fmt.Printf("  InvertDTR:      %t\n", p.InvertDTR != 0)     //
			fmt.Printf("  InvertDSR:      %t\n", p.InvertDSR != 0)     //
			fmt.Printf("  InvertDCD:      %t\n", p.InvertDCD != 0)     //
			fmt.Printf("  InvertRI:       %t\n", p.InvertRI != 0)      //
			fmt.Printf("  Cbus0:          %s\n", p.Cbus0)              //
			fmt.Printf("  Cbus1:          %s\n", p.Cbus1)              //
			fmt.Printf("  Cbus2:          %s\n", p.Cbus2)              //
			fmt.Printf("  Cbus3:          %s\n", p.Cbus3)              //
			fmt.Printf("  Cbus4:          %s\n", p.Cbus4)              //
			fmt.Printf("  DriverType:     %d\n", p.DriverType)         // 0 is D2XX, 1 is VCP
		default:
			fmt.Printf("Unknown device:   %s\n", i.Type)
		}
		log.Printf("  Raw: %x\n", ee.Raw)
	} else {
		fmt.Printf("Failed to read EEPROM: %v\n", err)
	}

	if ua, err := d.UserArea(); err != nil {
		fmt.Printf("Failed to read UserArea: %v\n", err)
	} else {
		fmt.Printf("UserArea: %x\n", ua)
	}

	hdr := d.Header()
	for _, p := range hdr {
		fmt.Printf("%s: %s\n", p, p.Function())
	}
}

func mainImpl() error {
	verbose := flag.Bool("v", false, "verbose mode")
	flag.Parse()
	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}
	log.SetFlags(log.Lmicroseconds)
	if flag.NArg() != 0 {
		return errors.New("unexpected argument, try -help")
	}

	if _, err := host.Init(); err != nil {
		return err
	}

	major, minor, build := d2xx.Version()
	fmt.Printf("Using library %d.%d.%d\n", major, minor, build)
	all := ftdi.All()
	plural := ""
	if len(all) > 1 {
		plural = "s"
	}
	fmt.Printf("Found %d device%s\n", len(all), plural)
	for i, d := range all {
		fmt.Printf("- Device #%d\n", i)
		process(d)
		if i != len(all)-1 {
			fmt.Printf("\n")
		}
	}
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "ftdi-list: %s.\n", err)
		os.Exit(1)
	}
}
