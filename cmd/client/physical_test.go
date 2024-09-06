package main_test

import (
	"github.com/channingko-madden/pi-vitrine/cmd/client"
	"testing"
)

func TestGetEnvData(t *testing.T) {
	env, err := main.GetEnvData()

	if err != nil {
		t.Errorf("Error measuring environment data %s", err)
	}

	t.Logf("%8s %10s %9s\n", env.Temperature, env.Pressure, env.Humidity)

}
