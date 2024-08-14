package system_test

import (
	"github.com/channingko-madden/pi-vitrine/internal/system"
	"testing"
)

func TestMeasureCPUTemp(t *testing.T) {
	temp, error := system.MeasureCPUTemp()
	if error != nil || temp == -273.15 {
		t.Errorf("Error measuring CPU temp %s", error)
	}
}
