package main

import (
	"github.com/channingko-madden/pi-vitrine/db"
	"log"
	"net/http"
)

// the _ in the import sets the name of the package to _.
// This is b/c this driver package shouldn't be used directly!

// DB is a handle, not the actual connection. It has a pool of
// DB connections in the background.
// Can use a global struct like this or pass around DB.
var Db *db.PostgresDeviceRepository

// the init function is called automatically for every package!
// Sets up connection to the database (doesn't open it!)
// Opening occurs lazily
func init() {
	connection := "user=pi-vitrine dbname=pi-vitrine password=pi-vitrine"

	Db = db.NewPostgresDeviceRepository(connection)
}

func main() {
	http.HandleFunc("/", HomePageHandler)
	http.HandleFunc("/system", CreateSystemDataHandler)
	log.Fatal(http.ListenAndServe(":9000", nil))
}
