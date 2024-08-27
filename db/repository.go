package db

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib" // self registers a postgres driver
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
