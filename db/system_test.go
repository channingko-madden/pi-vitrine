package db_test

import (
	"github.com/channingko-madden/pi-vitrine/db"
	"github.com/channingko-madden/pi-vitrine/db/dbtest"
	"github.com/channingko-madden/pi-vitrine/internal/cher"
	"testing"
	"time"
)

func TestSystemDataCreate(t *testing.T) {

	testDb := db.NewPostgresDeviceRepository(dbtest.CreateTestDbConnection(t))

	// Device must be created in order to create system data
	deviceData := cher.Device{
		Name: "my_name",
	}

	err := testDb.CreateDevice(&deviceData)
	if err != nil {
		t.Fatal(err)
	}

	systemData := cher.System{
		Name:    deviceData.Name,
		CPUTemp: 40.0,
		GPUTemp: 35.5,
	}

	err = testDb.CreateSystem(&systemData)
	if err != nil {
		t.Fatal(err)
	}

	if systemData.CreatedAt.IsZero() {
		t.Fail()
	}

	t.Logf("%+v\n", systemData)

	getData, err := testDb.GetSystemData(deviceData.Name, time.Time{}, time.Time{})

	if err != nil {
		t.Fatal(err)
	}

	if len(getData) != 1 {
		t.Fail()
	}

	if getData[0] != systemData {
		t.Fail()
	}

}
