package cher

type System struct {
	MacAddr   string  `json:"mac_addr"`
	CPUTemp   float64 `json:"cpu_temp"`
	GPUTemp   float64 `json:"gpu_temp"`
	CreatedAt string  `json:"created_at"`
}
