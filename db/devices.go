package db

import (
	"github.com/channingko-madden/pi-vitrine/internal/cher"
)

// Represent device information
// Devices are identified by their unique MAC address
// Should serial be used instead?
type Device struct {
	Id        int
	MacAddr   string
	Hardware  string
	Revision  string
	Serial    string
	Model     string
	CreatedAt string
}

type DeviceRepository interface {
	CreateDevice(device *cher.Device) error
	GetDevice(macAddr string) (cher.Device, error)
}

func (r *PostgresDeviceRepository) CreateDevice(deviceData *cher.Device) error {

	statement := "insert into devices (mac_addr, hardware, revision, serial, model) values ($1, $2, $3, $4, $5) returning created_at"
	stmt, err := r.conn.Prepare(statement)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(deviceData.MacAddr, deviceData.Hardware, deviceData.Revision, deviceData.Serial, deviceData.Model).Scan(&deviceData.CreatedAt)
	return err
}

func (r *PostgresDeviceRepository) GetDevice(deviceName string) (cher.Device, error) {

	statement := "select id, mac_addr, hardware, revision, serial, model, created_at from devices where name = $1"

	stmt, err := r.conn.Prepare(statement)
	if err != nil {
		return cher.Device{}, err
	}
	defer stmt.Close()

	var device Device
	err = stmt.QueryRow(deviceName).Scan(&device.Id, &device.MacAddr, &device.Hardware, &device.Revision, &device.Serial, &device.Model, &device.CreatedAt)

	if err != nil {
		return cher.Device{}, err
	}

	outDevice := cher.Device{
		Name:      deviceName,
		MacAddr:   device.MacAddr,
		Hardware:  device.Hardware,
		Revision:  device.Revision,
		Serial:    device.Serial,
		Model:     device.Model,
		CreatedAt: device.CreatedAt,
	}

	return outDevice, err
}
