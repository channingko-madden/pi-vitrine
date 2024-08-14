package db_test

import (
	"context"
	"database/sql"
	"github.com/channingko-madden/pi-vitrine/db"
	"github.com/channingko-madden/pi-vitrine/internal/data"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewDb(connection string) *sql.DB {
	Db, err := sql.Open("pgx", connection)

	if err != nil {
		panic(err)
	}
	return Db
}

func TestSystemDataCreate(t *testing.T) {

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

	testDb := NewDb(connStr)

	systemData := data.System{
		MacAddr: "macaddr",
		CPUTemp: 40.0,
		GPUTemp: 35.5,
	}

	err = db.CreateSystem(testDb, &systemData)

	if systemData.Id == 0 {
		t.Fail()
	}

	if len(systemData.CreatedAt) == 0 {
		t.Fail()
	}

	t.Logf("%+v\n", systemData)

	getData, err := db.GetAllSystemData(testDb, "macaddr")

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
