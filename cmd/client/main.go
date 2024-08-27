package main

import (
	_ "github.com/jackc/pgx/v5/stdlib" // self registers a postgres driver
	"log"
	"net/http"
)

// Pass URL of server as a CLI arg
func main() {
	http.HandleFunc("/", HomePageHandler)
	log.Fatal(http.ListenAndServe(":9000", nil))
}
