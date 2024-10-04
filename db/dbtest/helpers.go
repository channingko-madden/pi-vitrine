package dbtest

import (
	"github.com/channingko-madden/pi-vitrine/db"
	"github.com/channingko-madden/pi-vitrine/internal/cher"
	"math"
	"testing"
)

func CreateTestDevice(t *testing.T, r db.DeviceRepository, device_name string) {
	// Device must be created in order to add indoor climate data
	deviceData := cher.Device{
		Name: device_name,
	}

	err := r.CreateDevice(&deviceData)
	if err != nil {
		t.Fatal(err)
	}
}

// Return if the relative difference between two float64 values are smaller than a given epsilon
// value.
func CompareFloat(expected, value, epsilon float64) bool {
	if expected == value {
		return true
	}

	diff := math.Abs(expected - value)

	if expected == 0 {
		return diff < epsilon // avoid div by 0
	} else {
		return (diff / math.Abs(expected)) < epsilon
	}

}
