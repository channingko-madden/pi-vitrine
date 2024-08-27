package main

import (
	"fmt"
	"github.com/channingko-madden/pi-vitrine/db"
	"log"
	"net/http"
)

// DB is a handle, not the actual connection. It has a pool of
// DB connections in the background.
// Can use a global struct like this or pass around DB.
var Db *db.PostgresDeviceRepository

// the init function is called automatically for every package!
// Sets up connection to the database (doesn't open it!)
// Opening occurs lazily
func init() {
	connection := "user=pi-vitrine dbname=pi_vitrine password=pi-vitrine"

	Db = db.NewPostgresDeviceRepository(connection)
}

func main() {

	fmt.Println("localhost:9000")

	http.HandleFunc("GET /", HomePageHandler)
	http.HandleFunc("POST /system", CreateSystemDataHandler)
	http.HandleFunc("GET /system", GetSystemDataHandler)
	http.HandleFunc("GET /device", GetDeviceHandler)
	http.HandleFunc("POST /device", CreateDeviceHandler)
	log.Fatal(http.ListenAndServe(":9000", nil))
}
