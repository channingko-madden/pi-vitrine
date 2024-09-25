package main_test

import (
	"bytes"
	"encoding/json"
	"github.com/channingko-madden/pi-vitrine/cmd/server"
	"github.com/channingko-madden/pi-vitrine/internal"
	"github.com/channingko-madden/pi-vitrine/internal/cher"
	"net/http"
	"net/http/httptest"
	"testing"
)

type goodDb struct {
}

func (goodDb) CreateSystem(system *cher.System) error {
	return nil
}

func (goodDb) GetAllSystemData(macAddr string) ([]cher.System, error) {
	return []cher.System{
		{
			Name:    "test",
			GPUTemp: 12.3,
			CPUTemp: 45.6,
		},
	}, nil
}

func (goodDb) CreateDevice(device *cher.Device) error {
	return nil
}

func (goodDb) GetDevice(macAddr string) (cher.Device, error) {
	return cher.Device{Name: "test"}, nil
}
func (goodDb) GetAllDevices() ([]cher.Device, error) {
	return []cher.Device{
		{
			Name: "test",
		},
	}, nil
}

func TestCreateSystemDataHandler(t *testing.T) {

	main.Db = goodDb{}

	body := cher.System{
		Name:    "123ABC",
		CPUTemp: 23.32,
		GPUTemp: 45.45,
	}

	bytes_json, err := json.Marshal(body)

	if err != nil {
		t.Fatalf("Error marshalling SystemData: %s", err)
	}

	req := httptest.NewRequest("POST", "/system", bytes.NewReader(bytes_json))

	recorder := httptest.NewRecorder()

	handler := internal.HostErrorHandler(main.CreateSystemDataHandler)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("handler returned the wrong status code: got %v wanted %v",
			status, http.StatusCreated)
	}
}
