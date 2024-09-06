package main

import (
	"fmt"
	"log"
	"net/http"
)

// Pass URL of server as a CLI arg
func main() {
	fmt.Println("Running")

	// start blink
	go BlinkLED()

	http.HandleFunc("GET /", HomePageHandler)
	http.HandleFunc("GET /env", GetEnvHandler)
	log.Fatal(http.ListenAndServe(":9000", nil))
}
