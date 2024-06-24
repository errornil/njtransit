package main

import (
	"log"
	"net/http"
	"os"
	"time"

	njt "github.com/errornil/njtransit/v2"
)

func main() {
	log.Println("Starting...")

	hc := &http.Client{
		Timeout: 10 * time.Second,
	}

	client, err := njt.NewBusClient(
		njt.BusProdURL,
		os.Getenv("BUS_USERNAME"),
		os.Getenv("BUS_PASSWORD"),
		os.Getenv("USER_AGENT"),
		hc,
	)
	if err != nil {
		log.Fatalf("Failed to create BusClient: %v", err)
	}

	b, err := client.GetGTFS()
	if err != nil {
		log.Fatalf("Failed to call GetGTFS: %v", err)
	}

	// save to file
	err = os.WriteFile("gtfs.zip", b, 0644)
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	// feed, err := client.GetTripUpdates()
	// if err != nil {
	// 	log.Fatalf("Failed to call GetTripUpdates: %v", err)
	// }

	// for _, entity := range feed.GetEntity() {
	// 	tripUpdate := entity.GetTripUpdate()
	// 	log.Printf("Trip: %s", tripUpdate.GetTrip().GetTripId())
	// 	for _, stopTimeUpdate := range tripUpdate.GetStopTimeUpdate() {
	// 		log.Printf("  Stop: %s", stopTimeUpdate.GetStopId())
	// 	}
	// }
}
