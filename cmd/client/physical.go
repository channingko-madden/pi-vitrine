package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/channingko-madden/pi-vitrine/internal"
	"github.com/channingko-madden/pi-vitrine/internal/cher"
	"github.com/channingko-madden/pi-vitrine/internal/system"

	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/bmxx80"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
)

func init() {
	// init needs to be outside the goroutine
	if _, err := host.Init(); err != nil {
		panic(err)
	}
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
// Return Env data in units °C, kPa and % of relative humidity
func GetEnvData() (physic.Env, error) {
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
	var env physic.Env
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

			data.CPUTemp = internal.CelciusToKelvin(data.CPUTemp)

			data.GPUTemp, err = system.MeasureGPUTemp()

			if err != nil {
				log.Printf("Error getting GPU Temp: %s", err)
				continue
			}

			data.GPUTemp = internal.CelciusToKelvin(data.GPUTemp)

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
					body, _ := io.ReadAll(response.Body)
					log.Print("pi-vitrine client sent a system data POST that was not accepted: ", body)
					return
				}
			}
			<-t.C

		}
	}
}

// Send indoor climate data to the pi-vitrine server every 30 minutes
func SendIndoorClimateData(clientName string, serverAddress string, ctx context.Context) {

	t := time.NewTicker(30 * time.Minute)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			envData, err := GetEnvData()

			if err != nil {
				log.Print("Error reading indoor climate data: ", err)
			}

			// convert to the units the server expects! K, Pa, %
			climateData := cher.IndoorClimate{
				Name:             clientName,
				AirTemp:          float64(envData.Temperature) / float64(physic.Kelvin),
				Pressure:         float64(envData.Pressure) / float64(float64(physic.Pascal)),
				RelativeHumidity: float64(envData.Humidity) / float64(physic.PercentRH),
			}

			jsonData, err := json.Marshal(climateData)

			if err != nil {
				log.Print("Error json marshalling: ", err)
				continue
			}

			request, err := http.NewRequest("POST", serverAddress, bytes.NewBuffer(jsonData))

			if err != nil {
				log.Print("Error creating HTTP POST request: ", err)
				continue
			}

			request.Header.Set("Content-Type", "application/json; charset=UTF-8")

			client := http.Client{}
			response, err := client.Do(request)

			if err != nil {
				log.Print("Error sending HTTP POST request: ", err)
				continue
			}
			defer response.Body.Close()

			if response.StatusCode != 201 {
				log.Printf("pi-vitrine server returned an error code %d", response.StatusCode)
				if response.StatusCode == 400 {
					body, _ := io.ReadAll(response.Body)
					log.Print("pi-vitrine client sent an indoor_climate POST that was not accepted: ", body)
					return
				}
			}
			<-t.C
		}
	}
}
