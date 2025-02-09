package main_test

import (
	"github.com/channingko-madden/pi-vitrine/cmd/client"
	"periph.io/x/conn/v3/physic"
	"testing"
)

func TestGetEnvData(t *testing.T) {
	env, err := main.GetEnvData()

	if err != nil {
		t.Errorf("Error measuring environment data %s", err)
	}

	t.Logf("%8s %10s %9s\n", env.Temperature, env.Pressure, env.Humidity)


	temp :=          float64(env.Temperature) / float64(physic.Kelvin)
	pressure:=         float64(env.Pressure) / float64(float64(physic.Pascal))
	rh:= float64(env.Humidity) / float64(physic.PercentRH)

	t.Logf("%f (K) %f (P) %f (%%)\n", temp, pressure, rh)
}
