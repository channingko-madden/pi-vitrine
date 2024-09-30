package cher

type System struct {
	Name      string  `json:"device_name"`
	CPUTemp   float64 `json:"cpu_temp"` // Units are Kelvin
	GPUTemp   float64 `json:"gpu_temp"` // Units are Kelvin
	CreatedAt string
}
