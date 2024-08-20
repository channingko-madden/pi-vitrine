package main_test

import (
	"bytes"
	"encoding/json"
	"github.com/channingko-madden/pi-vitrine/cmd/server"
	"github.com/channingko-madden/pi-vitrine/db"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateSystemDataHandler(t *testing.T) {
	body := db.SystemData{
		MacAddr: "123ABC",
		CPUTemp: 23.32,
		GPUTemp: 45.45,
	}

	bytes_json, err := json.Marshal(body)

	if err != nil {
		t.Fatalf("Error marshalling SystemData: %s", err)
	}

	req := httptest.NewRequest("POST", "/system", bytes.NewReader(bytes_json))

	recorder := httptest.NewRecorder()

	handler := http.HandlerFunc(main.CreateSystemDataHandler)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("handler returned the wrong status code: got %v wanted %v",
			status, http.StatusCreated)
	}
}
