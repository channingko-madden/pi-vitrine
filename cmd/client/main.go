package main

import (
	"log"
	"net/http"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
	"time"
)

func blink_led() {
	t := time.NewTicker(500 * time.Millisecond)
	for l := gpio.Low; ; l = !l {
		rpi.P1_36.Out(l) // physical pin my boy!
		<-t.C
	}
}

// Pass URL of server as a CLI arg
func main() {

	host.Init() // init needs to be outside the goroutine
	// start blink
	go blink_led()

	http.HandleFunc("/", HomePageHandler)
	log.Fatal(http.ListenAndServe(":9000", nil))
}
