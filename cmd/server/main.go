package main

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib" // self registers a postgres driver
	"log"
	"net/http"
)

// the _ in the import sets the name of the package to _.
// This is b/c this driver package shouldn't be used directly!

// DB is a handle, not the actual connection. It has a pool of
// DB connections in the background.
// Can use a global struct like this or pass around DB.
var Db *sql.DB

// the init function is called automatically for every package!
// Sets up connection to the database (doesn't open it!)
// Opening occurs lazily
func init() {
	var err error
	// can't use := b/c scope?
	Db, err = sql.Open("pgx", "user=pi-vitrine dbname=pi-vitrine password=pi-vitrine")

	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", HomePageHandler)
	http.HandleFunc("/system", CreateSystemDataHandler)
	log.Fatal(http.ListenAndServe(":9000", nil))
}
