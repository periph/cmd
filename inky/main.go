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
	dcPin := flag.String("dc", "GPIO22", "Inky DC Pin")
	resetPin := flag.String("reset", "GPIO27", "Inky Reset Pin")
	busyPin := flag.String("busy", "GPIO17", "Inky Busy Pin")
	model := inky.IMPRESSION73
	flag.Var(&model, "model", "Inky model (PHAT, PHAT2, WHAT, IMPRESSION4, IMPRESSION57, or IMPRESSION73)")
	modelColor := inky.Multi
	flag.Var(&modelColor, "model-color", "Inky model color (multi, black, red or yellow)")
	borderColor := inky.Red
	flag.Var(&borderColor, "border-color", "Border color (multi, black, white, red or yellow)")
	flag.Parse()

	// Open and decode the image.
	f, err := os.Open(*path)
	if err != nil {
		return err
	}
	/* #nosec G307 */
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
	opts := &inky.Opts{
		Model:       model,
		ModelColor:  modelColor,
		BorderColor: borderColor,
	}

	log.Printf("Drawing image...")

	if model <= inky.PHAT2 {
		dev, eInky := inky.New(b, dc, reset, busy, opts)
		if eInky != nil {
			return eInky
		}
		return dev.Draw(img.Bounds(), img, image.Point{})
	}

	dev, err := inky.NewImpression(b, dc, reset, busy, opts)
	if err != nil {
		return err
	}
	return dev.Draw(img.Bounds(), img, image.Point{})
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "inky: %s.\n", err)
		os.Exit(1)
	}
}
