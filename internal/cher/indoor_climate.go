package cher

import "time"

type IndoorClimate struct {
	Name             string  `json:"device_name"`
	AirTemp          float64 `json:"air_temp"`          // Units are Kelvin
	Pressure         float64 `json:"pressure"`          // Units are Pascal
	RelativeHumidity float64 `json:"relative_humidity"` // A 0 - 100% value
	CreatedAt        string
}

func (i IndoorClimate) CreatedTime() time.Time {
	time, err := time.Parse(time.RFC3339Nano, i.CreatedAt)
	if err != nil {
		panic(err)
	}

	return time
}
