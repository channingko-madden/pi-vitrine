package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/channingko-madden/pi-vitrine/db"
	"github.com/channingko-madden/pi-vitrine/internal"
	"github.com/channingko-madden/pi-vitrine/internal/cher"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFS(content, "templates/home.html")
	if err != nil {
		panic(err)
	}
	temp.Execute(w, nil)
}

// post "/system"
// TODO: Create test for this, "testing" package has it all!
func CreateSystemDataHandler(w http.ResponseWriter, r *http.Request) *internal.HostError {
	decoder := json.NewDecoder(r.Body)

	var data cher.System
	err := decoder.Decode(&data)

	if err != nil {
		return &internal.HostError{
			Error:   err,
			Message: "Could not decode json body",
			Code:    400,
		}
	}

	err = Db.CreateSystem(&data)

	if err != nil {
		return &internal.HostError{
			Error:   err,
			Message: "Could not store system data",
			Code:    500,
		}
	} else {
		w.WriteHeader(201)
		return nil
	}
}

// get "/system" returns system data for a given device
// Returns html with data summaries

// TODO: query params to limit time range of data
func GetSystemDataHandler(w http.ResponseWriter, r *http.Request) {
	deviceName := r.URL.Query().Get("device_name")

	data, err := Db.GetAllSystemData(deviceName)

	if err != nil || len(data) == 0 {
		internal.ErrorMessage(w, fmt.Sprintf("System data for device %s not found", deviceName))
	} else {
		temp, err := template.ParseFS(content, "templates/system_data.html")
		if err != nil {
			panic(err)
		}
		temp.Execute(w, data[len(data)-1])
	}

}

// Responds with html
func GetDeviceHandler(w http.ResponseWriter, r *http.Request) {
	deviceName := r.PathValue("name")

	if len(deviceName) == 0 {
		internal.ErrorMessage(w, "URL path is missing 'name'")
		return
	}

	device, err := Db.GetDevice(deviceName)

	if err != nil {
		if dneError, ok := err.(*db.DeviceDoesNotExistError); ok {
			internal.ErrorMessage(w, dneError.Error())
		} else {
			internal.ErrorMessage(w, fmt.Sprintf("Server error finding device %s", deviceName))
		}
	} else {
		temp, err := template.ParseFS(content, "templates/device.html")
		if err != nil {
			panic(err)
		}
		temp.Execute(w, device)
	}
}

// Responds with html
func GetAllDevicesHandler(w http.ResponseWriter, r *http.Request) {

	temp, err := template.ParseFS(content, "templates/list_devices.html")
	if err != nil {
		panic(err)
	}

	devices, err := Db.GetAllDevices()
	if err != nil {
		if _, ok := err.(*db.DeviceDoesNotExistError); ok {
			internal.ErrorMessage(w, "No devices exist")
		} else {
			internal.ErrorMessage(w, "Server error listing devices")
		}
		return
	}

	temp.Execute(w, devices)
}

// Responds with html
func CreateDeviceHandler(w http.ResponseWriter, r *http.Request) {
	// deviceName is within the POST form!
	r.ParseForm() // for urlencoded!

	deviceName := r.PostFormValue("device_name")

	if len(deviceName) == 0 {
		internal.ErrorMessage(w, "POST form is missing 'device_name'")
		return
	}

	// check if device already exists
	existingDevice, err := Db.GetDevice(deviceName)
	if err != nil {
		if _, ok := err.(*db.DeviceDoesNotExistError); !ok {
			log.Default().Print(err)
			internal.ErrorMessage(w, fmt.Sprintf("Error saving device %s", deviceName))
			return
		}
	} else if existingDevice.Name == deviceName {
		internal.ErrorMessage(w, fmt.Sprintf("Device %s already exists", deviceName))
		return
	}

	device := cher.Device{
		Name: deviceName,
	}
	err = Db.CreateDevice(&device)

	if err != nil {
		log.Default().Print(err)
		internal.ErrorMessage(w, fmt.Sprintf("Error saving device %s", deviceName))
	} else {
		internal.InfoMessage(w, fmt.Sprintf("Created device %s", deviceName))
	}

}
