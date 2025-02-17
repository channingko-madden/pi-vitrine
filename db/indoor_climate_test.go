package db_test

import (
	"github.com/channingko-madden/pi-vitrine/db"
	"github.com/channingko-madden/pi-vitrine/db/dbtest"
	"github.com/channingko-madden/pi-vitrine/internal/cher"
	"testing"
	"time"
)

func TestCreateIndoorClimate(t *testing.T) {
	testDb := db.NewPostgresDeviceRepository(dbtest.CreateTestDbConnection(t))

	// Device must be created in order to add indoor climate data
	dbtest.CreateTestDevice(t, testDb, "my_name")

	climateData := cher.IndoorClimate{
		Name:             "my_name",
		AirTemp:          300.5,
		Pressure:         101325.5,
		RelativeHumidity: 50.5,
	}

	err := testDb.CreateIndoorClimate(&climateData)

	if err != nil {
		t.Fatal(err)
	}

	if climateData.CreatedAt.IsZero() {
		t.Error("CreatedAt missing from creating climate data")
	}

	allData, err := testDb.GetIndoorClimateData("my_name", time.Time{}, time.Time{})

	if err != nil {
		t.Fatal(err)
	}

	if len(allData) != 1 {
		t.Errorf("Expected 1 row of indoor_climate data, got %d", len(allData))
	}

	if !compareIndoorClimate(allData[0], climateData) {
		t.Log(climateData)
		t.Log(allData[0])
		t.Errorf("Expected indoor_climate data did not match returned data")
	}
}

func TestGetIndoorClimateData(t *testing.T) {
	testDb := db.NewPostgresDeviceRepository(dbtest.CreateTestDbConnection(t))

	testDeviceName := "my_name"

	// Device must be created in order to add indoor climate data
	dbtest.CreateTestDevice(t, testDb, testDeviceName)

	expectedData := []cher.IndoorClimate{
		{
			Name:             testDeviceName,
			AirTemp:          300.5,
			Pressure:         101325.5,
			RelativeHumidity: 50.5,
		},
		{
			Name:             testDeviceName,
			AirTemp:          302.5,
			Pressure:         101310.5,
			RelativeHumidity: 52.5,
		},
		{
			Name:             testDeviceName,
			AirTemp:          305.8,
			Pressure:         101300.25,
			RelativeHumidity: 53.25,
		},
		{
			Name:             testDeviceName,
			AirTemp:          306.34,
			Pressure:         101290.45,
			RelativeHumidity: 53.88,
		},
	}

	for i, expected := range expectedData {
		err := testDb.CreateIndoorClimate(&expected)
		if err != nil {
			t.Fatal(err)
		}
		expectedData[i] = expected
	}

	midIndex := len(expectedData) / 2
	t.Log("midIndex: ", midIndex)

	startTime := expectedData[0].CreatedAt
	midTime := expectedData[midIndex].CreatedAt
	endTime := expectedData[len(expectedData)-1].CreatedAt

	t.Log("startTime: ", startTime.Format(time.RFC3339Nano))
	t.Log("midTime: ", midTime.Format(time.RFC3339Nano))
	t.Log("endTime: ", endTime.Format(time.RFC3339Nano))

	// start = zero value, end = zero value
	allData, err := testDb.GetIndoorClimateData(testDeviceName, time.Time{}, time.Time{})

	if err != nil {
		t.Fatal(err)
	}

	if !compareIndoorClimateData(expectedData, allData) {
		t.Errorf("Passing zero value for start/end should return all Indoor Climate data")
		t.Log("Expected: ", expectedData)
		t.Log("Returned: ", allData)
	}

	// start = start, end = end
	allData, err = testDb.GetIndoorClimateData(testDeviceName, startTime, endTime)

	if err != nil {
		t.Fatal(err)
	}

	if !compareIndoorClimateData(expectedData, allData) {
		t.Errorf("Start & Ends should have return all Indoor Climate data")
		t.Log("Expected: ", expectedData)
		t.Log("Returned: ", allData)
	}

	// start = start, end = mid
	allData, err = testDb.GetIndoorClimateData(testDeviceName, startTime, midTime)

	if err != nil {
		t.Fatal(err)
	}

	subslice := expectedData[0 : midIndex+1]

	if !compareIndoorClimateData(subslice, allData) {
		t.Errorf("Start & Midpoint time as end did not return the expected data")
		t.Log("Expected: ", subslice)
		t.Log("Returned: ", allData)
	}

	// start = zero value, end = mid
	allData, err = testDb.GetIndoorClimateData(testDeviceName, time.Time{}, midTime)

	if err != nil {
		t.Fatal(err)
	}

	subslice = expectedData[0 : midIndex+1]

	if !compareIndoorClimateData(subslice, allData) {
		t.Errorf("Zero value start & Midpoint time as end did not return the expected data")
		t.Log("Expected: ", subslice)
		t.Log("Returned: ", allData)
	}

	// start = mid, end = end
	allData, err = testDb.GetIndoorClimateData(testDeviceName, midTime, endTime)

	if err != nil {
		t.Fatal(err)
	}

	subslice = expectedData[midIndex:]

	if !compareIndoorClimateData(subslice, allData) {
		t.Errorf("Midpoint time as start, and endTime did not return the expected data")
		t.Log("Expected: ", subslice)
		t.Log("Returned: ", allData)
	}

	// start = mid, end = zero value
	allData, err = testDb.GetIndoorClimateData(testDeviceName, midTime, time.Time{})

	if err != nil {
		t.Fatal(err)
	}

	if !compareIndoorClimateData(subslice, allData) {
		t.Errorf("Midpoint time as start, and zero value end did not return the expected data")
		t.Log("Expected: ", subslice)
		t.Log("Returned: ", allData)
	}
}

func compareIndoorClimateData(expected, value []cher.IndoorClimate) bool {
	if len(expected) != len(value) {
		return false
	}

	for i := range expected {
		if !compareIndoorClimate(expected[i], value[i]) {
			return false
		}
	}

	return true
}

func compareIndoorClimate(expected, value cher.IndoorClimate) bool {

	epsilon := 0.001

	return expected.Name == value.Name &&
		dbtest.CompareFloat(expected.AirTemp, value.AirTemp, epsilon) &&
		dbtest.CompareFloat(expected.Pressure, value.Pressure, epsilon) &&
		dbtest.CompareFloat(expected.RelativeHumidity, value.RelativeHumidity, epsilon) &&
		expected.CreatedAt == value.CreatedAt
}
