package db

import (
	"github.com/channingko-madden/pi-vitrine/internal/cher"
	"time"
)

type indoorClimate struct {
	Id               int
	DeviceId         int
	AirTemp          float64 // Kelivin
	Pressure         float64 // Pascal
	RelativeHumidity float64 // Decimal percent 0-100%
	CreatedAt        string  // Unix timestamp "YYYY-MM-DD HH::MM::SS"
}

type IndoorClimateRepository interface {
	CreateIndoorClimate(climate *cher.IndoorClimate) error

	// Return all indoor atmosphere data for a given device, within the reporting period [start, end].
	// Pass a time.Time zero value if start and/or end are not desired.
	GetIndoorClimateData(deviceName string, start time.Time, end time.Time) ([]cher.IndoorClimate, error)
}

func (r *PostgresDeviceRepository) CreateIndoorClimate(climate *cher.IndoorClimate) error {

	// Get the device id using the device name
	deviceId, err := r.getDeviceId(climate.Name)

	if err != nil {
		return err
	}

	statement := "insert into indoor_climate (device_id, air_temp, pressure, relative_humidity) values ($1, $2, $3, $4) returning created_at"
	stmt, err := r.conn.Prepare(statement)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(deviceId, climate.AirTemp, climate.Pressure, climate.RelativeHumidity).Scan(&climate.CreatedAt)
	return err
}

func (r *PostgresDeviceRepository) GetIndoorClimateData(deviceName string, start time.Time, end time.Time) ([]cher.IndoorClimate, error) {
	// Get the device id using the device name
	deviceId, err := r.getDeviceId(deviceName)

	if err != nil {
		return nil, err
	}

	statement := "select air_temp, pressure, relative_humidity, created_at from indoor_climate where device_id = $1"

	timeFilter := r.ReportingPeriodWhereFilter(start, end)
	if len(timeFilter) != 0 {
		statement = statement + " and " + timeFilter
	}

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

	var allData []cher.IndoorClimate
	for rows.Next() {
		climateData := indoorClimate{}
		err = rows.Scan(&climateData.AirTemp, &climateData.Pressure, &climateData.RelativeHumidity, &climateData.CreatedAt)
		if err != nil {
			return nil, err
		}

		outData := cher.IndoorClimate{
			Name:             deviceName,
			AirTemp:          climateData.AirTemp,
			Pressure:         climateData.Pressure,
			RelativeHumidity: climateData.RelativeHumidity,
			CreatedAt:        climateData.CreatedAt,
		}
		allData = append(allData, outData)
	}
	return allData, nil
}
