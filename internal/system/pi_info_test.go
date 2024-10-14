package system_test

import (
	"github.com/channingko-madden/pi-vitrine/internal/system"
	"testing"
)

func TestMeasureCPUTemp(t *testing.T) {
	temp, error := system.MeasureCPUTemp()
	if error != nil {
		t.Errorf("Error measuring CPU temp %s", error)
	}

	if temp <= -273.15 {
		t.Error("Measured CPU temp is below absolute 0...")
	}

	if temp > 85.0 {
		t.Logf("Measured CPU temp of %f is above the max temp for a raspberry pi!", temp)
	}
}

func TestMeasureGPUTemp(t *testing.T) {
	temp, error := system.MeasureGPUTemp()
	if error != nil {
		t.Errorf("Error measuring GPU temp %s", error)
	}

	if temp <= -273.15 {
		t.Error("Measured GPU temp is below absolute 0...")
	}

	if temp > 85.0 {
		t.Logf("Measured GPU temp of %f is above the max temp for a raspberry pi!", temp)
	}
}
