package data

type System struct {
	MacAddr   string  `json:"macaddr"`
	CPUTemp   float64 `json:"cputemp"`
	GPUTemp   float64 `json:"gputemp"`
	CreatedAt string
}
