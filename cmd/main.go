package main

import (
	"fmt"
	"github.com/channingko-madden/pi-vitrine/internal/system"
)

func main() {

	tempC, err := system.MeasureTemp()
	if err == nil {
		fmt.Println(tempC)
	}

	distro, err := system.DistroInfo()
	if err == nil {
		fmt.Println(distro)
	}

	pinfo := system.GetPiInfo()
	fmt.Println(pinfo)
}
