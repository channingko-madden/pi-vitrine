package db

import (
	"database/sql"
)

// Pi system information
type SystemData struct {
	Id        int
	MacAddr   string
	CPUTemp   float64
	GPUTemp   float64
	CreatedAt string
}

func (data *SystemData) Create(db *sql.DB) error {
	statement := "insert into system (mac_addr, temp_cpu, temp_gpu) values ($1, $2, $3) returning id, created_at"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(data.MacAddr, data.CPUTemp, data.GPUTemp).Scan(&data.Id, &data.CreatedAt)
	return err
}

func GetAllSystemData(db *sql.DB, macAddr string) ([]SystemData, error) {

	data := []SystemData{}

	statement := "select id, temp_cpu, temp_gpu, created_at from system where mac_addr = $1"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return data, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(macAddr)
	if err != nil {
		return data, err
	}
	defer rows.Close()
	for rows.Next() {
		systemData := SystemData{
			MacAddr: macAddr,
		}
		err = rows.Scan(&systemData.Id, &systemData.CPUTemp, &systemData.GPUTemp, &systemData.CreatedAt)
		if err != nil {
			return data, err
		}
		data = append(data, systemData)
	}
	return data, err
}
