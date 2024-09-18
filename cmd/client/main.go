package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Pass URL of server as a CLI arg
func main() {

	var clientAddressFlag = flag.String("client_address", "localhost", "IP address for the client homepage")
	var clientPortFlag = flag.Int("client_port", 9000, "Port number for the client homepage")

	flag.Parse()

	addr := *clientAddressFlag + ":" + strconv.Itoa(*clientPortFlag)

	fmt.Printf("pi-vitrine client running on %s\n", addr)

	// start blink
	go BlinkLED()

	http.HandleFunc("GET /", HomePageHandler)
	http.HandleFunc("GET /env", GetEnvHandler)
	log.Fatal(http.ListenAndServe(addr, nil))
}
