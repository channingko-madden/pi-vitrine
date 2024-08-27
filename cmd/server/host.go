package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/channingko-madden/pi-vitrine/internal/cher"
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

	var data cher.System
	err := decoder.Decode(&data)

	if err != nil {
		w.WriteHeader(400)
		return
	}

	err = Db.CreateSystem(&data)

	if err != nil {
		w.WriteHeader(500)
	} else {
		w.WriteHeader(201)
	}
}

// get "/system" returns system data for a given device
// query param "mac_addr" contains the MAC address
// Returns html with data summaries

// TODO: query params to limit time range of data
func GetSystemDataHandler(w http.ResponseWriter, r *http.Request) {
	macAddr := r.URL.Query().Get("mac_addr")

	data, err := Db.GetAllSystemData(macAddr)

	if err != nil || len(data) == 0 {
		errorMessage(w, fmt.Sprintf("System data for device with MAC Address %s not found", macAddr))
	} else {
		temp, err := template.ParseFiles("cmd/server/templates/system_data.html")
		if err == nil {
			temp.Execute(w, data[len(data)-1])
		} else {
			log.Default().Print(err)
		}
	}

}

func errorMessage(w http.ResponseWriter, errorMsg string) {
	temp, err := template.ParseFiles("cmd/server/templates/error_msg.html")
	if err == nil {
		temp.Execute(w, errorMsg)
	} else {
		log.Default().Print(err)
	}

}

// query param "mac_addr" contains the MAC address
// Returns html
func GetDeviceHandler(w http.ResponseWriter, r *http.Request) {
	macAddr := r.URL.Query().Get("mac_addr")

	device, err := Db.GetDevice(macAddr)

	if err != nil {
		errorMessage(w, fmt.Sprintf("Device with MAC Address %s not found", macAddr))
	} else {
		temp, err := template.ParseFiles("cmd/server/templates/device.html")
		if err == nil {
			temp.Execute(w, device)
		} else {
			log.Default().Print(err)
		}
	}
}

func CreateDeviceHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var data cher.Device
	err := decoder.Decode(&data)

	if err != nil {
		log.Default().Print(err)
		w.WriteHeader(400)
		return
	}

	err = Db.CreateDevice(&data)

	if err != nil {
		log.Default().Print(err)
		w.WriteHeader(500)
	} else {
		w.WriteHeader(201)
	}
}
