// Copyright 2021 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// videosink is a simplistic example of how to use the "videosink" virtual
// display as an HTTP handler.
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"net/http"
	"os"
	"time"

	"periph.io/x/devices/v3/videosink"
)

func main() {
	listenAddr := flag.String("listen-addr", ":8080", "Address and port to listen on")
	flag.Parse()

	dev := videosink.New(&videosink.Options{
		Width:  640,
		Height: 480,
	})
	defer func() {
		if err := dev.Halt(); err != nil {
			log.Printf("Halting dev failed: %v", err)
		}
	}()

	go func() {
		img := image.NewRGBA(dev.Bounds())
		draw.Draw(img, img.Bounds(), image.White, image.Point{}, draw.Src)

		size := int(img.Bounds().Dy() / 30)
		pos := 0

		fgcolor := &color.Gray{}
		fgcolor.Y = 0

		fg := &image.Uniform{C: fgcolor}

		for {
			draw.Draw(img, image.Rect(img.Bounds().Min.X, pos-(size/2), img.Bounds().Max.X, pos+(size/2)), fg, image.Point{}, draw.Src)

			pos += size

			if pos > img.Bounds().Dy() {
				pos = 0
				fgcolor.Y += 0xFF / 20
			}

			if err := dev.Draw(dev.Bounds(), img, image.Point{}); err != nil {
				log.Printf("Drawing failed: %v", err)
			}

			time.Sleep(time.Second / 30)
		}
	}()

	fmt.Fprintf(os.Stderr, "Starting HTTP server on %q\n", *listenAddr)

	http.Handle("/display", dev)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html><head></head>
			<body style="background: #ccc;">
			<img src="/display" width="%d" height="%d" style="border: 1px solid #000;">
			</body>
			</html>`,
			dev.Bounds().Dx(), dev.Bounds().Dy())
	})
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
