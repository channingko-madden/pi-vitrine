package db_test

import (
	"github.com/channingko-madden/pi-vitrine/db"
	"github.com/channingko-madden/pi-vitrine/db/dbtest"
	"github.com/channingko-madden/pi-vitrine/internal/cher"
	"reflect"
	"testing"
	"time"
)

func TestCreateIndoorClimate(t *testing.T) {
	testDb := db.NewPostgresDeviceRepository(dbtest.CreateTestDbConnection(t))

	// Device must be created in order to add indoor climate data
	deviceData := cher.Device{
		Name: "my_name",
	}

	err := testDb.CreateDevice(&deviceData)
	if err != nil {
		t.Fatal(err)
	}

	climateData := cher.IndoorClimate{
		Name:             "my_name",
		AirTemp:          300.5,
		Pressure:         101325.5,
		RelativeHumidity: 50.5,
	}

	err = testDb.CreateIndoorClimate(&climateData)

	if err != nil {
		t.Fatal(err)
	}

	if len(climateData.CreatedAt) == 0 {
		t.Error("CreatedAt missing from creating climate data")
	}

	allData, err := testDb.GetIndoorClimateData("my_name", time.Time{}, time.Time{})

	if err != nil {
		t.Fatal(err)
	}

	if len(allData) != 1 {
		t.Errorf("Expected 1 row of indoor_climate data, got %d", len(allData))
	}

	if !reflect.DeepEqual(allData[0], climateData) {
		t.Log(climateData)
		t.Log(allData[0])
		t.Errorf("Expected indoor_climate data did not match returned data")
	}
}
