package cher

// Represent device information
type Device struct {
	Name      string `json:"name"`
	MacAddr   string `json:"mac_addr"`
	Hardware  string `json:"hardware"`
	Revision  string `json:"revision"`
	Serial    string `json:"serial"`
	Model     string `json:"model"`
	CreatedAt string
}
