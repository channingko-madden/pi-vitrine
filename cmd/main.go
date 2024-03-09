package main

import (
	"github.com/channingko-madden/pi-vitrine/cmd/site"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", site.HomePageHandler)
	log.Fatal(http.ListenAndServe(":9000", nil))
}
