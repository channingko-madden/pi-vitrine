package db

import (
	"github.com/channingko-madden/pi-vitrine/internal/cher"
	"time"
)

// Pi system information
type systemData struct {
	Id        int
	DeviceId  int
	CPUTemp   float64
	GPUTemp   float64
	CreatedAt string
}

func newSystemData(s *cher.System) systemData {
	return systemData{
		CPUTemp: s.CPUTemp,
		GPUTemp: s.GPUTemp,
	}
}

type SystemRepository interface {
	CreateSystem(system *cher.System) error
	// Return all system data for a given device, within the reporting period [start, end].
	// Pass a time.Time zero value if start and/or end are not desired.
	GetSystemData(deviceName string, start time.Time, end time.Time) ([]cher.System, error)
}

func (r *PostgresDeviceRepository) CreateSystem(data *cher.System) error {

	// Get the device id using the device name
	deviceId, err := r.getDeviceId(data.Name)

	if err != nil {
		return err
	}

	statement := "insert into system (device_id, temp_cpu, temp_gpu) values ($1, $2, $3) returning created_at"
	stmt, err := r.Conn.Prepare(statement)
	if err != nil {
		return err
	}
	defer stmt.Close()

	s := newSystemData(data)

	err = stmt.QueryRow(deviceId, s.CPUTemp, s.GPUTemp).Scan(&s.CreatedAt)

	if err != nil {
		return nil
	}

	data.CreatedAt = parseCreatedTime(s.CreatedAt)

	return err
}

func (r *PostgresDeviceRepository) GetSystemData(deviceName string, start time.Time, end time.Time) ([]cher.System, error) {

	// Get the device id using the device name
	deviceId, err := r.getDeviceId(deviceName)

	if err != nil {
		return nil, err
	}

	statement := "select temp_cpu, temp_gpu, created_at from system where device_id = $1"

	timeFilter := r.ReportingPeriodWhereFilter(start, end)
	if len(timeFilter) != 0 {
		statement += " and " + timeFilter
	}

	stmt, err := r.Conn.Prepare(statement)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(deviceId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allData []cher.System
	for rows.Next() {
		systemData := systemData{}
		err = rows.Scan(&systemData.CPUTemp, &systemData.GPUTemp, &systemData.CreatedAt)
		if err != nil {
			return nil, err
		}

		outData := cher.System{
			Name:      deviceName,
			CPUTemp:   systemData.CPUTemp,
			GPUTemp:   systemData.GPUTemp,
			CreatedAt: parseCreatedTime(systemData.CreatedAt),
		}
		allData = append(allData, outData)
	}
	return allData, nil
}
