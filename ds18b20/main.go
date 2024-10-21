// Copyright 2024 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.
//
// This program enumerates onewire buses for  DS18B20 sensors and
// continuously reading them.
//
// If you get the error no buses were found, ensure that onewire buses
// are enabled. On a Raspberry Pi, run raspi-config, Interface Options,
// and enable 1-Wire. Remember to reboot after enabling it.
package main

import (
	"fmt"
	"log"
	"time"

	"periph.io/x/conn/v3/onewire"
	"periph.io/x/conn/v3/onewire/onewirereg"
	"periph.io/x/devices/v3/ds18b20"
	"periph.io/x/host/v3"
)

// For the Dallas onewire devices, the conversion time is dependent on the
// resolution. Refer to the Datasheet for more information.
const DefaultBits = 10

// enumerateBusForSensors searches for addresses on the bus and if it finds
// a DS18X20 series device, returns it.
func enumerateBusForSensors(bus onewire.Bus) []*ds18b20.Dev {
	result := make([]*ds18b20.Dev, 0)
	addresses, err := bus.Search(false)
	if err != nil {
		log.Print("  Search error:", err)
		return result
	}

	for _, address := range addresses {
		family := ds18b20.Family(address & 0xff)
		if family == ds18b20.DS18S20 || family == ds18b20.DS18B20 {
			log.Printf("Found Device %s Address %#016X\n", family, address)
			dev, err := ds18b20.New(bus, address, DefaultBits)
			if err != nil {
				log.Print(err)
			} else {
				result = append(result, dev)
			}
		}
	}
	return result
}

func main() {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	buses := make([]onewire.Bus, 0)
	sensors := make([]*ds18b20.Dev, 0)
	for _, ref := range onewirereg.All() {
		log.Printf("Found One-Wire Bus %s\n", ref.Name)
		busCloser, err := ref.Open()
		if err != nil {
			fmt.Println(" Open error:", err)
			continue
		}
		bus := busCloser.(onewire.Bus)
		newSensors := enumerateBusForSensors(bus)
		if len(newSensors) > 0 {
			// Do an initial convert on all devices on the bus.
			err = ds18b20.ConvertAll(bus, DefaultBits)
			if err != nil {
				log.Println(err)
			}
			buses = append(buses, bus)
			sensors = append(sensors, newSensors...)
		}
	}
	if len(buses) == 0 {
		log.Fatal("no onewire buses found.")
	}
	if len(sensors) == 0 {
		log.Fatal("no DS18X20 sensors found.")
	}

	// This demo uses StartAll()/LastTemp() to read without any latency.
	// You can also use Sense() to read the units, but there will be
	// the conversion latency.
	for {
		out := fmt.Sprintf("%d ", time.Now().Unix())
		for _, dev := range sensors {
			t, err := dev.LastTemp()
			if err == nil {
				out += fmt.Sprintf("%.2f ", t.Celsius())
			} else {
				log.Print(err)
				out += "err "
			}

		}
		fmt.Println(out)

		for _, bus := range buses {
			// Start the conversion cycle running while we're sleeping ...
			err := ds18b20.StartAll(bus)
			if err != nil {
				log.Println(err)
			}
		}
		time.Sleep(time.Duration(1000-(time.Now().UnixMilli()%1000)) * time.Millisecond)
	}
}
