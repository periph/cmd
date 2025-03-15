// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// mpu9250 calibrates and performs the self-test, then measures the acceleration continuously.
package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"

	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"

	"periph.io/x/devices/v3/mpu9250"
	"periph.io/x/devices/v3/mpu9250/accelerometer"
)

var (
	ifType      = flag.String("iftype", "spi", "Interface Type (spi, i2c)")
	i2cAddr     = flag.Int("i2caddr", 0x68, "I2C Address - Default 0x68")
	accRes      = flag.String("accRes", "2", "Acceleration resolution (2, 4, 8, 16G)")
	continuous  = flag.Bool("cont", false, "Continuous read")
	sensitivity int
)

func main() {

	flag.Parse()

	switch *accRes {
	case "2", "2G", "2g":
		sensitivity = accelerometer.ACCEL_FS_SEL_2G
	case "4", "4G", "4g":
		sensitivity = accelerometer.ACCEL_FS_SEL_4G
	case "8", "8G", "8g":
		sensitivity = accelerometer.ACCEL_FS_SEL_8G
	case "16", "16G", "16g":
		sensitivity = accelerometer.ACCEL_FS_SEL_16G
	default:
		sensitivity = accelerometer.ACCEL_FS_SEL_2G
	}

	if _, err := host.Init(); err != nil {
		log.Fatal("Error initializing host", err)
	}

	var t *mpu9250.Transport
	var err error

	if *ifType == "spi" {
		cs := gpioreg.ByName("8")
		if cs == nil {
			log.Fatal("Can't initialize CS pin")
		}
		t, err = mpu9250.NewSpiTransport("", cs)
	} else {
		var bus i2c.Bus
		bus, err = i2creg.Open("")
		if err == nil {
			t, err = mpu9250.NewI2cTransport(bus, uint16(*i2cAddr))
		}
	}
	if err != nil {
		log.Fatalf("Can't initialize %s bus: %s", *ifType, err)
	}

	dev, err := mpu9250.New(*t)
	if err != nil {
		log.Fatal(err)
	}

	if err = dev.Init(); err != nil {
		log.Fatal(err)
	}

	dev.Debug(func(msg string, args ...interface{}) {
		fmt.Printf(msg, args...)
	})

	id, err := dev.GetDeviceID()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Dev ID: %x\n", id)

	st, err := dev.SelfTest()
	if err != nil {
		log.Fatal("Self test failed: ", err)
	}

	if err = dev.Calibrate(); err != nil {
		log.Fatal("Can't calibrate: ", err)
	}

	fmt.Printf("Accelerometer Deviation: X: %.2f%%, Y: %.2f%%, Z:%.2f%%\n", st.AccelDeviation.X, st.AccelDeviation.Y, st.AccelDeviation.Z)
	fmt.Printf("Gyroscope Deviation: X: %.2f%%, Y: %.2f%%, Z:%.2f%%\n", st.GyroDeviation.X, st.GyroDeviation.Y, st.GyroDeviation.Z)

	time.Sleep(40 * time.Millisecond)

	accMulti := accelerometer.Sensitivity(sensitivity)

	if err := dev.SetAccelRange(byte(sensitivity)); err != nil {
		log.Fatal(err)
	}

	if *continuous {
		for {
			x := mustInt16(dev.GetAccelerationX())
			y := mustInt16(dev.GetAccelerationY())
			z := mustInt16(dev.GetAccelerationZ())
			fmt.Printf("Raw : X: %d, Y: %d, Z: %d\n", x, y, z)
			fmt.Printf("Calc: X: %.2f, Y: %.2f, Z: %.2f\n", float32(x)*accMulti, float32(y)*accMulti, float32(z)*accMulti)
			time.Sleep(time.Second)
			fmt.Println("----------------")
		}
	}
}

func mustInt16(s int16, err error) int16 {
	if err != nil {
		panic(err)
	}
	return s
}
