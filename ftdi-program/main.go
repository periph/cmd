// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// ftdi-eeprom interacts with the EEPROM of a FTDI device.
//
// It can either program the EEPROM or the User Area, or read it back.
package main

import (
	"encoding/hex"
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

func writeEEPROM(d ftdi.Dev, manufacturer, manufacturerID, desc, serial string) error {
	ee := ftdi.EEPROM{}
	if err := d.EEPROM(&ee); err != nil {
		fmt.Printf("Failed to read EEPROM: %v\n", err)
	}
	ee.Manufacturer = manufacturer
	ee.ManufacturerID = manufacturerID
	ee.Desc = desc
	ee.Serial = serial
	log.Printf("Writing: %x", ee.Raw)
	return d.WriteEEPROM(&ee)
}

func mainImpl() error {
	verbose := flag.Bool("v", false, "verbose mode")
	erase := flag.Bool("e", false, "erases the EEPROM instead of programming it")
	manufacturer := flag.String("m", "", "manufacturer")
	manufacturerID := flag.String("mid", "", "manufacturer ID")
	desc := flag.String("d", "", "description")
	serial := flag.String("s", "", "serial")
	ua := flag.String("ua", "", "hex encoded data")

	flag.Parse()
	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}
	log.SetFlags(log.Lmicroseconds)

	if flag.NArg() != 0 {
		return errors.New("unexpected argument, try -help")
	}
	if *erase {
		if *ua != "" || *manufacturer != "" || *manufacturerID != "" || *desc != "" || *serial != "" {
			return errors.New("-e cannot be used with any of -m, -mid, -d, -s, -ua")
		}
	} else {
		if *ua == "" {
			if *manufacturer == "" || *manufacturerID == "" || *desc == "" || *serial == "" {
				return errors.New("all of -m, -mid, -d and -s are required, or use -ua")
			}
		} else {
			if *manufacturer != "" || *manufacturerID != "" || *desc != "" || *serial != "" {
				return errors.New("all of -m, -mid, -d and -s cannot be used with -ua")
			}
		}
	}

	if _, err := host.Init(); err != nil {
		return err
	}
	major, minor, build := d2xx.Version()
	log.Printf("Using library %d.%d.%d\n", major, minor, build)

	all := ftdi.All()
	if len(all) == 0 {
		return errors.New("found no FTDI device on the USB bus")
	}
	if len(all) > 1 {
		return fmt.Errorf("for safety reasons, plug exactly one FTDI device on the USB bus, found %d devices", len(all))
	}
	d := all[0]

	if *erase {
		return d.EraseEEPROM()
	}
	if *ua == "" {
		return writeEEPROM(d, *manufacturer, *manufacturerID, *desc, *serial)
	}
	raw, err := hex.DecodeString(*ua)
	if err != nil {
		return err
	}
	return d.WriteUserArea(raw)
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "ftdi-eeprom: %s.\n", err)
		os.Exit(1)
	}
}
