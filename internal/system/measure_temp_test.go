package system

import (
	"testing"
)

func TestMeasureCPUTemp(t *testing.T) {
	temp, error := MeasureCPUTemp()
	if error != nil || temp == -273.15 {
		t.Errorf("Error measuring CPU temp %s", error)
	}
}
