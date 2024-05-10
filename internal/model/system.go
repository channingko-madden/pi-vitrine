package model

import (
	"database/sql"
	"time"
)

type System struct {
	Time    time.Time
	TempCPU float64
}

// Save System data into db
func (system *System) Save(db *sql.DB) {
	// generate sql query
	statement := "insert into system (time, temp_cpu) values ($1 $2)"

}
