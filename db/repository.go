package db

import (
	"database/sql"
)

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
