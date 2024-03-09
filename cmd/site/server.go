package site

import (
	"html/template"
	"log"
	"net/http"

	"github.com/channingko-madden/pi-vitrine/internal/system"
)

type HomePageInfo struct {
	PiInfo     system.PiInfo
	Temp       float64
	DistroInfo string
}

func gatherHomePageData() HomePageInfo {
	var info = HomePageInfo{}

	info.Temp, _ = system.MeasureTemp()
	info.DistroInfo, _ = system.DistroInfo()
	info.PiInfo = system.GetPiInfo()

	return info
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("cmd/site/templates/home.html")
	if err == nil {
		data := gatherHomePageData()
		temp.Execute(w, data)
	} else {
		log.Default().Print(err)
	}
}
