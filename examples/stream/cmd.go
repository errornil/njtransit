package main

import (
	"fmt"
	"log"
	"os"
	"time"

	njt "github.com/errornil/njtransit"
)

func main() {
	log.Println("Starting...")

	client := njt.NewBusDataClient(
		os.Getenv("BUSDATA_USERNAME"),
		os.Getenv("BUSDATA_PASSWORD"),
		njt.BusDataProdURL,
	)

	r := make(chan njt.BusVehicleDataRow)
	e := make(chan error)

	go client.GetBusVehicleDataStream(r, e, 5*time.Second, true)

	log.Println("Listening stream...")
	for {
		select {
		case row := <-r:
			fmt.Println(row)
		case err := <-e:
			fmt.Println(err)
		}
	}
}
