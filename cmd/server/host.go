package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/channingko-madden/pi-vitrine/db"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("cmd/server/templates/home.html")
	if err == nil {
		temp.Execute(w, nil)
	} else {
		log.Default().Print(err)
	}
}

// post "/system"
// TODO: Create test for this, "testing" package has it all!
func CreateSystemDataHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var data db.SystemData
	err := decoder.Decode(&data)

	if err != nil {
		w.WriteHeader(400)
		return
	}

	err = data.Create(Db)

	if err != nil {
		w.WriteHeader(500)
	} else {
		w.WriteHeader(201)
	}
}

// get "/system/{mac_address}" returns system data for a given device
// query params to limiting time range of data
