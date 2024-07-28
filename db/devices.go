package db

import (
	"database/sql"
)

// Represent device information
// Devices are identified by their unique MAC address
// Should serial be used instead?
type Device struct {
	Id       int
	MacAddr  string
	Hardware string
	Revision string
	Serial   string
	Model    string
}

func (device *Device) Create(db *sql.DB) error {
	statement := "insert into devices (mac_addr, hardware, revision, serial, model) values ($1, $2, $3, $4, $5) returning id"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(device.MacAddr, device.Hardware, device.Revision, device.Serial, device.Model).Scan(&device.Id)
	return err
}
