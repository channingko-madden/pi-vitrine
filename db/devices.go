package db

import (
	"github.com/channingko-madden/pi-vitrine/internal/data"
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

type DeviceRepository interface {
	CreateDevice(device *data.Device) error
}

func (r *PostgresDeviceRepository) CreateDevice(deviceData *data.Device) error {

	device := Device{
		MacAddr:  deviceData.MacAddr,
		Hardware: deviceData.Hardware,
		Revision: deviceData.Revision,
		Serial:   deviceData.Serial,
		Model:    deviceData.Model,
	}

	statement := "insert into devices (mac_addr, hardware, revision, serial, model) values ($1, $2, $3, $4, $5) returning id"
	stmt, err := r.conn.Prepare(statement)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(device.MacAddr, device.Hardware, device.Revision, device.Serial, device.Model).Scan(&device.Id)
	return err
}
