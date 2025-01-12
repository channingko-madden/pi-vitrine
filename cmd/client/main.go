package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
)

func main() {

	var clientAddressFlag = flag.String("client_address", "localhost:9000", "Network address for the client homepage")

	var clientNameFlag = flag.String("name", "", "Name of the client")

	var serverAddressFlag = flag.String("server_address", "", "IP address of the pi-vitrine server")

	flag.Parse()

	// Flag validation
	if len(*clientNameFlag) == 0 {
		log.Fatal("Must provide the name of the client")
	}

	if len(*serverAddressFlag) == 0 {
		log.Fatal("Must provide the pi-vitrine server address")
	} else if _, _, err := net.SplitHostPort(*serverAddressFlag); err != nil {
		log.Fatalf("Must provide a valid network address for the pi-vitrine server\n%s", err)
	}

	if _, _, err := net.SplitHostPort(*clientAddressFlag); err != nil {
		log.Fatalf("Must provide a valid network address for the client\n%s", err)
	}

	fmt.Printf("pi-vitrine client running on %s\n", *clientAddressFlag)

	// One context to cancel all goroutines
	ctx, cancel := context.WithCancel(context.Background())
	// start blink
	go BlinkLED(ctx)

	// start sending system data
	go SendSystemData(*clientNameFlag, *serverAddressFlag+"/system", ctx)

	// start sending indoor climate data
	go SendIndoorClimateData(*clientNameFlag, *serverAddressFlag+"/indoor_climate", ctx)

	http.HandleFunc("GET /", HomePageHandler)
	http.HandleFunc("GET /env", GetEnvHandler)
	log.Fatal(http.ListenAndServe(*clientAddressFlag, nil))

	cancel() // stop goroutines
}
