package db_test

import (
	"context"
	"github.com/channingko-madden/pi-vitrine/db"
	"github.com/channingko-madden/pi-vitrine/internal/cher"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func CreateTestDbConnection(t *testing.T) string {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithInitScripts(filepath.Join("..", "testdata", "setup.sql")),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	return connStr
}

func TestSystemDataCreate(t *testing.T) {

	testDb := db.NewPostgresDeviceRepository(CreateTestDbConnection(t))

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

	if len(systemData.CreatedAt) == 0 {
		t.Fail()
	}

	t.Logf("%+v\n", systemData)

	getData, err := testDb.GetAllSystemData(deviceData.Name)

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

func TestDeviceCreate(t *testing.T) {
	testDb := db.NewPostgresDeviceRepository(CreateTestDbConnection(t))

	deviceData := cher.Device{
		Name: "my_name",
	}

	err := testDb.CreateDevice(&deviceData)
	if err != nil {
		t.Fatal(err)
	}

	if len(deviceData.CreatedAt) == 0 {
		t.Fatal(err)
	}

	t.Logf("%+v\n", deviceData)

	getData, err := testDb.GetDevice(deviceData.Name)

	if err != nil {
		t.Fatal(err)
	}

	if len(getData.CreatedAt) == 0 {
		t.Fatal(err)
	}

	if getData.CreatedAt != deviceData.CreatedAt {
		t.Fail()
	}

	if getData.Name != deviceData.Name {
		t.Fail()
	}

	t.Logf("%+v\n", getData)

}

func TestCreatingSameDevice(t *testing.T) {
	testDb := db.NewPostgresDeviceRepository(CreateTestDbConnection(t))

	deviceData := cher.Device{
		Name: "my_name",
	}

	err := testDb.CreateDevice(&deviceData)
	if err != nil {
		t.Fatal(err)
	}

	if len(deviceData.CreatedAt) == 0 {
		t.Fatal(err)
	}

	deviceData.CreatedAt = ""

	err = testDb.CreateDevice(&deviceData)

	if err == nil {
		t.Fatal(err)
	}

	t.Log(err)
}

func TestGettingAllDevices(t *testing.T) {
	testDb := db.NewPostgresDeviceRepository(CreateTestDbConnection(t))

	expectedDevices := []cher.Device{
		{
			Name: "my_name",
		},
		{
			Name: "my_other_name",
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
