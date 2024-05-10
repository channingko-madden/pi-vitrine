package system

import (
	"os/exec"
	"strconv"
	"strings"
)

// Measure pi temperature using /usr/bin/vcgencmd measure_temp
// This command returns a value in the format "temp=XX.X'C
func MeasureGPUTemp() (float64, error) {
	command := exec.Command("/usr/bin/vcgencmd", "measure_temp")
	out, err := command.Output()
	if err != nil {
		return -273.15, err
	}
	return strconv.ParseFloat(string(out[5:9]), 64)
}

func MeasureCPUTemp() (float64, error) {
	command := exec.Command("cat", "/sys/class/thermal/thermal_zone0/temp")
	out, err := command.Output()
	if err != nil {
		return -273.15, err
	}
	temp, err := strconv.ParseFloat(strings.TrimSuffix(string(out), "\n"), 64)
	if err != nil {
		return -273.15, err
	}
	return temp / 1000.0, nil
}

// Return the contents of /proc/version file
func DistroInfo() (string, error) {
	command := exec.Command("cat", "/proc/version")
	out, err := command.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// Return the contenst of /proc/cpuinfo file
func cpuInfo() ([]byte, error) {
	command := exec.Command("cat", "/proc/cpuinfo")
	return command.Output()
}

// Represent the information about the Pi Board
type PiInfo struct {
	Hardware string
	Revision string
	Serial   string
	Model    string
}

// extracted from the /proc/cpuinfo file
func GetPiInfo() PiInfo {
	info, err := cpuInfo()
	if err != nil {
		return PiInfo{}
	}

	// parse out data we are interested in
	lines := strings.Split(string(info), "\n")

	var piInfo = PiInfo{}

	for _, l := range lines {
		if strings.HasPrefix(l, "Hardware") {
			if _, after, found := strings.Cut(l, ": "); found {
				piInfo.Hardware = after
			}
		} else if strings.HasPrefix(l, "Revision") {
			if _, after, found := strings.Cut(l, ": "); found {
				piInfo.Revision = after
			}
		} else if strings.HasPrefix(l, "Serial") {
			if _, after, found := strings.Cut(l, ": "); found {
				piInfo.Serial = after
			}
		} else if strings.HasPrefix(l, "Model") {
			if _, after, found := strings.Cut(l, ": "); found {
				piInfo.Model = after
			}
		}
	}

	return piInfo
}
