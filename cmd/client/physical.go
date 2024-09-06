package main

import (
	"log"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/bmxx80"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
	"time"
)

func init() {
	host.Init() // init needs to be outside the goroutine
}

func BlinkLED() {
	t := time.NewTicker(500 * time.Millisecond)
	for l := gpio.Low; ; l = !l {
		rpi.P1_36.Out(l) // physical pin my boy!
		<-t.C
	}
}

// Assumes host has been Init prior to calling this function!
func GetEnvData() (physic.Env, error) {
	var env physic.Env

	if _, err := host.Init(); err != nil {
		log.Println(err)
		return env, err
	}

	bus, err := i2creg.Open("1")
	if err != nil {
		log.Println(err)
		return env, err
	}
	defer bus.Close()

	// Connects to BME280 on I2C bus using default settings
	dev, err := bmxx80.NewI2C(bus, 0x77, &bmxx80.DefaultOpts)
	if err != nil {
		log.Println(err)
		return env, err
	}
	defer dev.Halt()

	// Read data
	if err = dev.Sense(&env); err != nil {
		log.Println(err)
		return env, err
	}

	return env, nil
}
