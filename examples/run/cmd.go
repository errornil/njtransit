package main

import (
	"log"
	"os"

	njt "github.com/errornil/njtransit"
)

func main() {
	log.Println("Starting...")

	client := njt.NewBusDataClient(
		os.Getenv("BUSDATA_USERNAME"),
		os.Getenv("BUSDATA_PASSWORD"),
		njt.BusDataProdURL,
	)

	log.Println("Calling GetBusVehicleData...")
	resp, err := client.GetBusVehicleData()
	if err != nil {
		log.Fatalf("Failed to call GetBusDV: %v", err)
	}
	log.Printf("%#v", *resp)

	log.Println("Calling GetNextTrips...")
	resp2, err := client.GetNextTrips(njt.GetNextTripsRequest{StopID: 21820})
	if err != nil {
		log.Fatalf("Failed to call GetNextTrips: %v", err)
	}
	log.Printf("%#v", *resp2)

	log.Println("Calling GetBusDV...")
	resp3, err := client.GetBusDV(njt.GetBusDVRequest{Location: "PABT"})
	if err != nil {
		log.Fatalf("Failed to call GetBusDV: %v", err)
	}
	log.Printf("%#v", *resp3)

	log.Println("Calling GetBusLocations...")
	resp4, err := client.GetBusLocations()
	if err != nil {
		log.Fatalf("Failed to call GetBusLocations: %v", err)
	}
	log.Printf("%#v", *resp4)

	log.Println("Calling GetMessages...")
	resp5, err := client.GetMessages(njt.GetMessagesRequest{StopID: 21820})
	if err != nil {
		log.Fatalf("Failed to call GetMessages: %v", err)
	}
	log.Printf("%#v", *resp5)
}
