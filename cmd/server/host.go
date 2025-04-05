package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/channingko-madden/pi-vitrine/db"
	"github.com/channingko-madden/pi-vitrine/internal"
	"github.com/channingko-madden/pi-vitrine/internal/cher"
)

// GET "/"
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
// Returns html with data plots
// Path parameters:
//   - "device_name"
//
// Query parameters:
//   - "days" (The number of past days to get data for. Will be > 0 enforced by client, error if it
//     is not)

// TODO: query params to limit time range of data
func GetSystemDataHandler(w http.ResponseWriter, r *http.Request) {
	deviceName := r.PathValue("device_name")

	if len(deviceName) == 0 {
		internal.ErrorMessage(w, "URL path is missing 'device_name'")
		return
	}

	startTime, endTime, err := calcStartEndTime(w, r)

	if err != nil {
		return
	}

	data, err := Db.GetSystemData(deviceName, startTime, endTime)

	if err != nil {
		if dneError, ok := err.(*db.DeviceDoesNotExistError); ok {
			internal.ErrorMessage(w, dneError.Error())
		} else {
			log.Print(err)
			internal.ErrorMessage(w, "Server error retrieving system data")
		}
		return
	}

	if len(data) == 0 {
		internal.InfoMessage(w, "There is no system data available")
		return
	}

	err = chartSystemData(w, data)
	if err != nil {
		internal.ErrorMessage(w, "Server error rendering graph")
		log.Print(err)
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
	// device_name is within the POST form!
	r.ParseForm() // for urlencoded!

	deviceName := r.PostFormValue("device_name")

	if len(deviceName) == 0 {
		internal.ErrorMessage(w, "POST form is missing 'device_name'")
		return
	}

	deviceLocation := r.PostFormValue("device_location")
	if len(deviceLocation) == 0 {
		internal.ErrorMessage(w, "POST form is missing 'device_location'")
		return
	}

	// check if device already exists
	existingDevice, err := Db.GetDevice(deviceName)
	if err != nil {
		if _, ok := err.(*db.DeviceDoesNotExistError); !ok {
			log.Print(err)
			internal.ErrorMessage(w, fmt.Sprintf("Error saving device %s", deviceName))
			return
		}
	} else if existingDevice.Name == deviceName {
		internal.ErrorMessage(w, fmt.Sprintf("Device %s already exists", deviceName))
		return
	}

	device := cher.Device{
		Name:     deviceName,
		Location: deviceLocation,
	}
	err = Db.CreateDevice(&device)

	if err != nil {
		log.Print(err)
		internal.ErrorMessage(w, fmt.Sprintf("Error saving device %s", deviceName))
	} else {
		internal.InfoMessage(w, fmt.Sprintf("Created device %s", deviceName))
	}
}

/*
PATCH /device/{name}
"device_location"
Responds with html
*/
func UpdateDeviceHandler(w http.ResponseWriter, r *http.Request) {
	deviceName := r.PathValue("name")

	if len(deviceName) == 0 {
		internal.ErrorMessage(w, "URL path is missing 'name'")
		return
	}

	r.ParseForm() // for urlencoded!
	deviceLocation := r.PostFormValue("device_location")
	if len(deviceLocation) == 0 {
		internal.ErrorMessage(w, "POST form is missing 'device_location'")
		return
	}

	device, err := Db.GetDevice(deviceName)

	if err != nil {
		if dneError, ok := err.(*db.DeviceDoesNotExistError); ok {
			internal.ErrorMessage(w, dneError.Error())
		} else {
			internal.ErrorMessage(w, fmt.Sprintf("Server error finding device %s", deviceName))
		}
		return
	}

	if device.Location != deviceLocation {
		device.Location = deviceLocation

		err = Db.UpdateDevice(&device)
		if err != nil {
			internal.ErrorMessage(w, fmt.Sprintf("Server error updating device %s", deviceName))
			return
		}
	}

	temp, err := template.ParseFS(content, "templates/device.html")
	if err != nil {
		panic(err)
	}
	temp.Execute(w, device)
}

// POST "indoor_climate"
// A non html endpoint
func CreateIndoorClimateDataHandler(w http.ResponseWriter, r *http.Request) *internal.HostError {
	decoder := json.NewDecoder(r.Body)

	var data cher.IndoorClimate
	err := decoder.Decode(&data)

	if err != nil {
		return &internal.HostError{
			Error:   err,
			Message: "Could not decode json body",
			Code:    400,
		}
	}

	err = Db.CreateIndoorClimate(&data)

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

// GET "indoor_climate"
// An html endpoint
// Path parameters:
//   - "device_name"
//
// Query parameters:
//   - "days" (The number of past days to get data for. Will be > 0 enforced by client, error if it
//     is not)
func GetIndoorClimateChartHandler(w http.ResponseWriter, r *http.Request) {
	deviceName := r.PathValue("device_name")

	if len(deviceName) == 0 {
		internal.ErrorMessage(w, "URL path is missing 'device_name'")
		return
	}

	startTime, endTime, err := calcStartEndTime(w, r)

	if err != nil {
		return
	}

	devices, err := Db.GetIndoorClimateData(deviceName, startTime, endTime)
	if err != nil {
		if dneError, ok := err.(*db.DeviceDoesNotExistError); ok {
			internal.ErrorMessage(w, dneError.Error())
		} else {
			log.Print(err)
			internal.ErrorMessage(w, "Server error retrieving climate data")
		}
		return
	}

	if len(devices) == 0 {
		internal.InfoMessage(w, "There is no indoor climate data available")
		return
	}

	err = chartIndoorClimate(w, devices)
	if err != nil {
		internal.ErrorMessage(w, "Server error rendering graph")
		log.Print(err)
	}
}

// Return start time and end time calculated from the "days" query parameter
// Writes error message and returns and error if "days" query parameter is missing
func calcStartEndTime(w http.ResponseWriter, r *http.Request) (time.Time, time.Time, error) {
	days, err := strconv.Atoi(r.URL.Query().Get("days"))

	if err != nil {
		internal.ErrorMessage(w, "URL query is missing an integer value for 'days'")
		return time.Time{}, time.Time{}, err
	}

	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -days)
	return startTime, endTime, nil
}
