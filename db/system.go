package db

import (
	"github.com/channingko-madden/pi-vitrine/internal/cher"
)

// Pi system information
type SystemData struct {
	Id        int
	DeviceId  int
	CPUTemp   float64
	GPUTemp   float64
	CreatedAt string
}

type SystemRepository interface {
	CreateSystem(system *cher.System) error
	GetAllSystemData(macAddr string) ([]cher.System, error)
}

func (r *PostgresDeviceRepository) CreateSystem(data *cher.System) error {

	// Get the device id using the device name
	deviceId, err := r.getDeviceId(data.Name)

	if err != nil {
		return err
	}

	statement := "insert into system (device_id, temp_cpu, temp_gpu) values ($1, $2, $3) returning created_at"
	stmt, err := r.conn.Prepare(statement)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(deviceId, data.CPUTemp, data.GPUTemp).Scan(&data.CreatedAt)
	return err
}

func (r *PostgresDeviceRepository) GetAllSystemData(deviceName string) ([]cher.System, error) {

	var allData []cher.System

	// Get the device id using the device name
	deviceId, err := r.getDeviceId(deviceName)

	if err != nil {
		return nil, err
	}

	statement := "select temp_cpu, temp_gpu, created_at from system where device_id = $1"
	stmt, err := r.conn.Prepare(statement)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(deviceId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		systemData := SystemData{}
		err = rows.Scan(&systemData.CPUTemp, &systemData.GPUTemp, &systemData.CreatedAt)
		if err != nil {
			return nil, err
		}

		outData := cher.System{
			Name:      deviceName,
			CPUTemp:   systemData.CPUTemp,
			GPUTemp:   systemData.GPUTemp,
			CreatedAt: systemData.CreatedAt,
		}
		allData = append(allData, outData)
	}
	return allData, nil
}
