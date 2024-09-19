package cher

type System struct {
	Name      string  `json:"device_name"`
	CPUTemp   float64 `json:"cpu_temp"`
	GPUTemp   float64 `json:"gpu_temp"`
	CreatedAt string  `json:"created_at"`
}
