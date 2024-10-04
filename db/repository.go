package db

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib" // self registers a postgres driver
	"time"
)

type Repository interface {
	SystemRepository
	DeviceRepository
	IndoorClimateRepository
}

type PostgresDeviceRepository struct {
	conn *sql.DB
}

func NewPostgresDeviceRepository(connection string) *PostgresDeviceRepository {

	Db, err := sql.Open("pgx", connection)

	if err != nil {
		panic(err)
	}

	return &PostgresDeviceRepository{conn: Db}

}

// Return a string containing sql for limiting records to a given time period [start, end]
//
// Pass a time.Time zero value if start and/or end are not desired.
//
// Assumes the column with timestamp data is named "created_at"
func (r *PostgresDeviceRepository) ReportingPeriodWhereFilter(start time.Time, end time.Time) string {
	if start.IsZero() && end.IsZero() {
		return ""
	} else if start.IsZero() && !end.IsZero() {
		// all records up to end
		return fmt.Sprintf("created_at <= '%s'", end.Format(time.RFC3339Nano))
	} else if !start.IsZero() && end.IsZero() {
		// all records from start
		return fmt.Sprintf("created_at >= '%s'", start.Format(time.RFC3339Nano))
	} else {
		return fmt.Sprintf("created_at >= '%s' and created_at <= '%s'", start.Format(time.RFC3339Nano), end.Format(time.RFC3339Nano))
	}
}
