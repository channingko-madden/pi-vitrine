package db_test

import (
	"github.com/channingko-madden/pi-vitrine/db"
	"github.com/channingko-madden/pi-vitrine/db/dbtest"
	"github.com/channingko-madden/pi-vitrine/internal/cher"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func TestDeviceCreate(t *testing.T) {
	testDb := db.NewPostgresDeviceRepository(dbtest.CreateTestDbConnection(t))

	deviceData := cher.Device{
		Name:     "my_name",
		Location: "my_location",
	}

	err := testDb.CreateDevice(&deviceData)
	if err != nil {
		t.Fatal(err)
	}

	if deviceData.CreatedAt.IsZero() {
		t.Fatal(err)
	}

	t.Logf("%+v\n", deviceData)

	getData, err := testDb.GetDevice(deviceData.Name)

	if err != nil {
		t.Fatal(err)
	}

	if getData.CreatedAt.IsZero() {
		t.Fatal(err)
	}

	if getData.CreatedAt != deviceData.CreatedAt {
		t.Fail()
	}

	if getData.Name != deviceData.Name {
		t.Fail()
	}

	if getData.Location != deviceData.Location {
		t.Fail()
	}

	t.Logf("%+v\n", getData)

}

func TestCreatingDeviceSameName(t *testing.T) {
	testDb := db.NewPostgresDeviceRepository(dbtest.CreateTestDbConnection(t))

	deviceData := cher.Device{
		Name:     "my_name",
		Location: "my_locations",
	}

	err := testDb.CreateDevice(&deviceData)
	if err != nil {
		t.Fatal(err)
	}

	if deviceData.CreatedAt.IsZero() {
		t.Fatal(err)
	}

	deviceData.Location = "new_location" // will still fail b/c Name is the same

	err = testDb.CreateDevice(&deviceData)

	if err == nil {
		t.Fatal(err)
	}

	t.Log(err)
}

func TestGettingAllDevices(t *testing.T) {
	testDb := db.NewPostgresDeviceRepository(dbtest.CreateTestDbConnection(t))

	expectedDevices := []cher.Device{
		{
			Name:     "my_name",
			Location: " my_location",
		},
		{
			Name:     "my_other_name",
			Location: " my_other_location",
		},
	}

	for i, expected := range expectedDevices {
		err := testDb.CreateDevice(&expected)
		if err != nil {
			t.Fatal(err)
		}
		expectedDevices[i] = expected
	}

	allDevices, err := testDb.GetAllDevices()
	if err != nil {
		t.Fatal(err)
	}

	if len(allDevices) != 2 {
		t.Errorf("Was not returned the expected 2 devices")
	}

	slices.SortFunc(expectedDevices, func(i, j cher.Device) int {
		return strings.Compare(i.Name, j.Name)
	})

	slices.SortFunc(allDevices, func(i, j cher.Device) int {
		return strings.Compare(i.Name, j.Name)
	})

	if !reflect.DeepEqual(expectedDevices, allDevices) {
		t.Log(expectedDevices)
		t.Log(allDevices)
		t.Errorf("Expected devices did not match return devices")
	}
}

func TestUpdateDevice(t *testing.T) {
	testDb := db.NewPostgresDeviceRepository(dbtest.CreateTestDbConnection(t))

	deviceData := cher.Device{
		Name:     "my_name",
		Location: "my_location",
	}

	err := testDb.CreateDevice(&deviceData)
	if err != nil {
		t.Fatal(err)
	}

	if deviceData.CreatedAt.IsZero() {
		t.Fatal(err)
	}

	deviceData.Location = "my_new_location"

	err = testDb.UpdateDevice(&deviceData)
	if err != nil {
		t.Fatal(err)
	}

	getData, err := testDb.GetDevice(deviceData.Name)

	if err != nil {
		t.Fatal(err)
	}

	if getData.Location != deviceData.Location {
		t.Fail()
	}
}
