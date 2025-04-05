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
	Location  string
	CreatedAt string
}

func newDevice(dev *cher.Device) device {
	return device{
		Name:     dev.Name,
		Location: dev.Location,
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

	// Returns a DeviceDoesNotExistError if the device cannot be found
	UpdateDevice(device *cher.Device) error

	// Returns a DeviceDoesNotExistError if the device cannot be found
	GetDevice(macAddr string) (cher.Device, error)

	// Returns a DeviceDoesNotExistError if there are no devices
	GetAllDevices() ([]cher.Device, error)
}

func (r *PostgresDeviceRepository) CreateDevice(deviceData *cher.Device) error {

	statement := "insert into devices (name, location) values ($1, $2) returning created_at"
	stmt, err := r.Conn.Prepare(statement)
	if err != nil {
		return err
	}
	defer stmt.Close()

	dev := newDevice(deviceData)
	err = stmt.QueryRow(dev.Name, dev.Location).Scan(&dev.CreatedAt)

	if err != nil {
		return err
	}

	deviceData.CreatedAt = parseCreatedTime(dev.CreatedAt)

	return nil
}

// Device name is not updated
func (r *PostgresDeviceRepository) UpdateDevice(deviceData *cher.Device) error {

	statement := "update devices set location = $1 where name = $2"
	stmt, err := r.Conn.Prepare(statement)
	if err != nil {
		return err
	}
	defer stmt.Close()

	dev := newDevice(deviceData)
	_, err = stmt.Exec(dev.Location, dev.Name)

	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresDeviceRepository) GetDevice(deviceName string) (cher.Device, error) {

	statement := "select id, created_at, location from devices where name = $1"

	stmt, err := r.Conn.Prepare(statement)
	if err != nil {
		return cher.Device{}, err
	}
	defer stmt.Close()

	var dev device
	err = stmt.QueryRow(deviceName).Scan(&dev.Id, &dev.CreatedAt, &dev.Location)

	if err != nil {
		if err == sql.ErrNoRows {
			return cher.Device{}, &DeviceDoesNotExistError{Name: deviceName}

		}
		return cher.Device{}, err
	}

	outDevice := cher.Device{
		Name:      deviceName,
		Location:  dev.Location,
		CreatedAt: parseCreatedTime(dev.CreatedAt),
	}

	return outDevice, nil
}

func (r *PostgresDeviceRepository) GetAllDevices() ([]cher.Device, error) {
	var devices []cher.Device

	statement := "select name, created_at, location from devices"

	stmt, err := r.Conn.Prepare(statement)
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
		if err := rows.Scan(&dev.Name, &dev.CreatedAt, &dev.Location); err != nil {
			return nil, err
		}

		devices = append(devices, cher.Device{Name: dev.Name, CreatedAt: parseCreatedTime(dev.CreatedAt), Location: dev.Location})
	}

	return devices, nil
}

func (r *PostgresDeviceRepository) getDeviceId(deviceName string) (int, error) {
	statement := "select id from devices where name = $1"

	stmt, err := r.Conn.Prepare(statement)
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
