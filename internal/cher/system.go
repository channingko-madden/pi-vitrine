package cher

type System struct {
	Name      string  `json:"name"`
	CPUTemp   float64 `json:"cpu_temp"`
	GPUTemp   float64 `json:"gpu_temp"`
	CreatedAt string  `json:"created_at"`
}
