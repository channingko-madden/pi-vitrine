package db

import (
	"database/sql"
	"github.com/channingko-madden/pi-vitrine/internal/data"
)

// Pi system information
type SystemData struct {
	Id        int
	MacAddr   string
	CPUTemp   float64
	GPUTemp   float64
	CreatedAt string
}

type SystemRepository interface {
	CreateSystem(system *data.System) error
	GetAllSystemData(macAddr string) ([]data.System, error)
}

func (r *PostgresDeviceRepository) CreateSystem(data *data.System) error {
	statement := "insert into system (mac_addr, temp_cpu, temp_gpu) values ($1, $2, $3) returning created_at"
	stmt, err := r.conn.Prepare(statement)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(data.MacAddr, data.CPUTemp, data.GPUTemp).Scan(&data.CreatedAt)
	return err
}

func (r *PostgresDeviceRepository) GetAllSystemData(macAddr string) ([]data.System, error) {

	allData := []data.System{}

	statement := "select id, temp_cpu, temp_gpu, created_at from system where mac_addr = $1"
	stmt, err := r.conn.Prepare(statement)
	if err != nil {
		return allData, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(macAddr)
	if err != nil {
		return allData, err
	}
	defer rows.Close()
	for rows.Next() {
		systemData := SystemData{
			MacAddr: macAddr,
		}
		err = rows.Scan(&systemData.Id, &systemData.CPUTemp, &systemData.GPUTemp, &systemData.CreatedAt)
		if err != nil {
			return allData, err
		}

		outData := data.System{
			MacAddr:   systemData.MacAddr,
			CPUTemp:   systemData.CPUTemp,
			GPUTemp:   systemData.GPUTemp,
			CreatedAt: systemData.CreatedAt,
		}
		allData = append(allData, outData)
	}
	return allData, err
}
