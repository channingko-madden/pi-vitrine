package db

import (
	"database/sql"
	"fmt"
)

type PostgresDeviceRepository struct {
	conn *sql.DB
}

func NewPostgresDeviceRepository(user string, dbname string, password string) *PostgresDeviceRepository {

	s := fmt.Sprintf("user=%s dbname=%s password=%s", user, dbname, password)

	Db, err := sql.Open("pgx", s)

	if err != nil {
		panic(err)
	}

	return &PostgresDeviceRepository{conn: Db}

}
