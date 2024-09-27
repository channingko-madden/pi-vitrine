package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/channingko-madden/pi-vitrine/internal/cher"
	"github.com/channingko-madden/pi-vitrine/internal/system"
	"log"
	"net/http"

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

func BlinkLED(ctx context.Context) {
	t := time.NewTicker(500 * time.Millisecond)
	for l := gpio.Low; ; l = !l {
		select {
		case <-ctx.Done():
			rpi.P1_36.Out(gpio.Low)
			return
		default:
			rpi.P1_36.Out(l) // physical pin my boy!
			<-t.C
		}
	}
}

// Assumes host has been Init prior to calling this function!
func GetEnvData() (physic.Env, error) {
	var env physic.Env

	if _, err := host.Init(); err != nil {
		log.Println(err)
		return physic.Env{}, err
	}

	bus, err := i2creg.Open("1")
	if err != nil {
		log.Println(err)
		return physic.Env{}, err
	}
	defer bus.Close()

	// Connects to BME280 on I2C bus using default settings
	dev, err := bmxx80.NewI2C(bus, 0x77, &bmxx80.DefaultOpts)
	if err != nil {
		log.Println(err)
		return physic.Env{}, err
	}
	defer dev.Halt()

	// Read data
	if err = dev.Sense(&env); err != nil {
		log.Println(err)
		return physic.Env{}, err
	}

	return env, nil
}

// Send system data to the pi-vitrine server every 30 minutes
func SendSystemData(clientName string, serverAddress string, ctx context.Context) {

	var data = cher.System{
		Name: clientName,
	}
	var err error
	t := time.NewTicker(30 * time.Minute)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			data.CPUTemp, err = system.MeasureCPUTemp()

			if err != nil {
				log.Printf("Error getting CPU Temp: %s", err)
				continue
			}

			data.GPUTemp, err = system.MeasureGPUTemp()

			if err != nil {
				log.Printf("Error getting GPU Temp: %s", err)
				continue
			}

			jsonData, err := json.Marshal(data)

			if err != nil {
				log.Printf("Error json marshalling: %s", err)
				continue
			}

			request, err := http.NewRequest("POST", serverAddress, bytes.NewBuffer(jsonData))

			if err != nil {
				log.Printf("Error creating HTTP POST request: %s", err)
				continue
			}

			request.Header.Set("Content-Type", "application/json; charset=UTF-8")

			client := http.Client{}
			response, err := client.Do(request)

			if err != nil {
				log.Printf("Error sending HTTP POST request: %s", err)
				continue
			}
			defer response.Body.Close()

			if response.StatusCode != 201 {
				log.Printf("pi-vitrine server returned an error code %d", response.StatusCode)
				if response.StatusCode == 400 {
					log.Print("pi-vitrine client is not POSTing valid system data json!")
					return
				}
			}
			<-t.C

		}
	}
}
