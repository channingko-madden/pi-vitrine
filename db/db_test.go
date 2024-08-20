package db_test

import (
	"context"
	"github.com/channingko-madden/pi-vitrine/db"
	"github.com/channingko-madden/pi-vitrine/internal/cher"
	"path/filepath"
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

	systemData := cher.System{
		MacAddr: "macaddr",
		CPUTemp: 40.0,
		GPUTemp: 35.5,
	}

	err := testDb.CreateSystem(&systemData)
	if err != nil {
		t.Fatal(err)
	}

	if len(systemData.CreatedAt) == 0 {
		t.Fail()
	}

	t.Logf("%+v\n", systemData)

	getData, err := testDb.GetAllSystemData("macaddr")

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
		MacAddr:  "my_address",
		Hardware: "blazingly_fast",
		Revision: "newest",
		Serial:   "123abc",
		Model:    "B+",
	}

	err := testDb.CreateDevice(&deviceData)
	if err != nil {
		t.Fatal(err)
	}

	if len(deviceData.CreatedAt) == 0 {
		t.Fatal(err)
	}

	t.Logf("%+v\n", deviceData)

	getData, err := testDb.GetDevice(deviceData.MacAddr)

	if err != nil {
		t.Fatal(err)
	}

	if len(getData.CreatedAt) == 0 {
		t.Fatal(err)
	}

	if getData.CreatedAt != deviceData.CreatedAt {
		t.Fail()
	}

	if getData.MacAddr != deviceData.MacAddr {
		t.Fail()
	}

	if getData.Hardware != deviceData.Hardware {
		t.Fail()
	}

	if getData.Revision != deviceData.Revision {
		t.Fail()
	}

	if getData.Serial != deviceData.Serial {
		t.Fail()
	}

	if getData.Model != deviceData.Model {
		t.Fail()
	}

	t.Logf("%+v\n", getData)

}
