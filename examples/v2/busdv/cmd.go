package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	njt "github.com/errornil/njtransit/v2"
)

func main() {
	log.Println("Starting...")

	hc := &http.Client{
		Timeout: 60 * time.Second,
	}

	client, err := njt.NewBusDV2Client(
		njt.BusDVProdURL,
		os.Getenv("BUS_USERNAME"),
		os.Getenv("BUS_PASSWORD"),
		os.Getenv("USER_AGENT"),
		hc,
	)
	if err != nil {
		log.Fatalf("Failed to create BusClient: %v", err)
	}

	// dv, err := client.GetBusDV("31198", )
	dv, err := client.GetVehicleLocations(
		"40.737169",
		"-74.169868",
		2000,
		"ALL",
	)
	if err != nil {
		log.Fatalf("Failed to call GetVehicleLocations: %v", err)
	}

	b, err := json.MarshalIndent(dv, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal feed: %v", err)
	}

	// print to stdout
	fmt.Println(string(b))
}
