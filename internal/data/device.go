package data

// Represent device information
type Device struct {
	MacAddr  string `json:"macaddr"`
	Hardware string `json:"hardware"`
	Revision string `json:"revision"`
	Serial   string `json:"serial"`
	Model    string `json:"model"`
}
