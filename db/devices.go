package db

import (
	"database/sql"
	"fmt"

	"github.com/channingko-madden/pi-vitrine/internal/cher"
)

// Represent device information
// Devices are identified by their unique name
type Device struct {
	Id        int
	Name      string
	CreatedAt string
}

type DeviceDoesNotExistError struct {
	Name string
}

func (e *DeviceDoesNotExistError) Error() string {
	return fmt.Sprintf("Device %s does not exist", e.Name)
}

type DeviceRepository interface {
	CreateDevice(device *cher.Device) error
	GetDevice(macAddr string) (cher.Device, error)
	GetAllDevices() ([]cher.Device, error)
}

func (r *PostgresDeviceRepository) CreateDevice(deviceData *cher.Device) error {

	statement := "insert into devices (name) values ($1) returning created_at"
	stmt, err := r.conn.Prepare(statement)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(deviceData.Name).Scan(&deviceData.CreatedAt)
	return err
}

func (r *PostgresDeviceRepository) GetDevice(deviceName string) (cher.Device, error) {

	statement := "select id, created_at from devices where name = $1"

	stmt, err := r.conn.Prepare(statement)
	if err != nil {
		return cher.Device{}, err
	}
	defer stmt.Close()

	var device Device
	err = stmt.QueryRow(deviceName).Scan(&device.Id, &device.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return cher.Device{}, &DeviceDoesNotExistError{Name: deviceName}

		}
		return cher.Device{}, err
	}

	outDevice := cher.Device{
		Name:      deviceName,
		CreatedAt: device.CreatedAt,
	}

	return outDevice, nil
}

func (r *PostgresDeviceRepository) GetAllDevices() ([]cher.Device, error) {
	var devices []cher.Device

	statement := "select name, created_at from devices"

	stmt, err := r.conn.Prepare(statement)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &DeviceDoesNotExistError{}
		}
		return nil, err
	}

	for rows.Next() {
		var device cher.Device
		if err := rows.Scan(&device.Name, &device.CreatedAt); err != nil {
			return nil, err
		}

		devices = append(devices, device)
	}

	return devices, nil
}

func (r *PostgresDeviceRepository) getDeviceId(deviceName string) (int, error) {
	statement := "select id from devices where name = $1"

	stmt, err := r.conn.Prepare(statement)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var deviceId int
	err = stmt.QueryRow(deviceName).Scan(&deviceId)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, &DeviceDoesNotExistError{Name: deviceName}
		}
		return 0, err
	}

	return deviceId, nil
}
