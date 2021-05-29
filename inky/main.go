// Copyright 2019 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// inky tests an inky board.
package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"

	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/devices/v3/inky"
	"periph.io/x/host/v3"
)

func mainImpl() error {
	spiPort := flag.String("spi", "SPI0.0", "Name or number of SPI port to open")
	path := flag.String("image", "", "Path to a png file to display on the inky")
	dcPin := flag.String("dc", "22", "Inky DC Pin")
	resetPin := flag.String("reset", "27", "Inky Reset Pin")
	busyPin := flag.String("busy", "17", "Inky Busy Pin")
	model := inky.PHAT
	flag.Var(&model, "model", "Inky model (PHAT or WHAT)")
	modelColor := inky.Red
	flag.Var(&modelColor, "model-color", "Inky model color (black, red or yellow)")
	borderColor := inky.Black
	flag.Var(&borderColor, "border-color", "Border color (black, white, red or yellow)")
	flag.Parse()

	// Open and decode the image.
	f, err := os.Open(*path)
	if err != nil {
		return err
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return err
	}

	if _, err = host.Init(); err != nil {
		return err
	}

	log.Printf("Opening %s...", *spiPort)
	b, err := spireg.Open(*spiPort)
	if err != nil {
		return err
	}

	log.Printf("Opening pins...")
	dc := gpioreg.ByName(*dcPin)
	if dc == nil {
		return fmt.Errorf("invalid DC pin name: %s", *dcPin)
	}

	reset := gpioreg.ByName(*resetPin)
	if reset == nil {
		return fmt.Errorf("invalid Reset pin name: %s", *resetPin)
	}

	busy := gpioreg.ByName(*busyPin)
	if busy == nil {
		return fmt.Errorf("invalid Busy pin name: %s", *busyPin)
	}

	log.Printf("Creating inky...")
	dev, err := inky.New(b, dc, reset, busy, &inky.Opts{
		Model:       model,
		ModelColor:  modelColor,
		BorderColor: borderColor,
	})
	if err != nil {
		return err
	}

	log.Printf("Drawing image...")
	return dev.Draw(img.Bounds(), img, image.Point{})
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "inky: %s.\n", err)
		os.Exit(1)
	}
}
