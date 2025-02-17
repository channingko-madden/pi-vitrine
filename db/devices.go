package db

import (
	"database/sql"
	"fmt"

	"github.com/channingko-madden/pi-vitrine/internal/cher"
)

// Represent device information
// Devices are identified by their unique name
type device struct {
	Id        int
	Name      string
	CreatedAt string
}

func newDevice(dev *cher.Device) device {
	return device{
		Name: dev.Name,
	}
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

	dev := newDevice(deviceData)
	err = stmt.QueryRow(dev.Name).Scan(&dev.CreatedAt)

	if err != nil {
		return err
	}

	deviceData.CreatedAt = parseCreatedTime(dev.CreatedAt)

	return nil
}

func (r *PostgresDeviceRepository) GetDevice(deviceName string) (cher.Device, error) {

	statement := "select id, created_at from devices where name = $1"

	stmt, err := r.conn.Prepare(statement)
	if err != nil {
		return cher.Device{}, err
	}
	defer stmt.Close()

	var dev device
	err = stmt.QueryRow(deviceName).Scan(&dev.Id, &dev.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return cher.Device{}, &DeviceDoesNotExistError{Name: deviceName}

		}
		return cher.Device{}, err
	}

	outDevice := cher.Device{
		Name:      deviceName,
		CreatedAt: parseCreatedTime(dev.CreatedAt),
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
		var dev device
		if err := rows.Scan(&dev.Name, &dev.CreatedAt); err != nil {
			return nil, err
		}

		devices = append(devices, cher.Device{Name: dev.Name, CreatedAt: parseCreatedTime(dev.CreatedAt)})
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
