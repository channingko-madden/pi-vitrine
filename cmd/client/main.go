package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// Pass URL of server as a CLI arg
func main() {

	var clientAddressFlag = flag.String("client_address", "localhost", "IP address for the client homepage")
	var clientPortFlag = flag.Int("client_port", 9000, "Port number for the client homepage")

	var clientNameFlag = flag.String("name", "", "Name of the client")

	var serverAddressFlag = flag.String("server_address", "", "IP address of the pi-vitrine server")

	flag.Parse()

	if len(*clientNameFlag) == 0 {
		log.Fatal("Must provide the name of the client")
	}

	if len(*serverAddressFlag) == 0 {
		log.Fatal("Must provide the pi-vitrine server address")
	} else if _, err := url.ParseRequestURI(*serverAddressFlag); err != nil {
		log.Fatalf("Must provide a valid URI for the pi-vitrine server\n%s", err)
	}

	addr := *clientAddressFlag + ":" + strconv.Itoa(*clientPortFlag)

	fmt.Printf("pi-vitrine client running on %s\n", addr)

	// One context to cancel all goroutines
	ctx, cancel := context.WithCancel(context.Background())
	// start blink
	go BlinkLED(ctx)

	// start sending system data
	go SendSystemData(*clientNameFlag, *serverAddressFlag, ctx)

	http.HandleFunc("GET /", HomePageHandler)
	http.HandleFunc("GET /env", GetEnvHandler)
	log.Fatal(http.ListenAndServe(addr, nil))

	cancel() // stop goroutines
}
