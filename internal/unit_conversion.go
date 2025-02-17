package internal

func CelciusToKelvin(celcius float64) float64 {
	return celcius + 273.15
}

func KelvinToCelcius(kelvin float64) float64 {
	return kelvin - 273.15
}
