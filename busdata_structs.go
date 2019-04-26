package njtransit

// GetBusVehicleDataResponse represents GetBusVehicleData API response
type GetBusVehicleDataResponse struct {
	Rows []BusVehicleDataRow `xml:"ROW"`
}

// BusVehicleDataRow is part of the GetBusVehicleDataResponse
type BusVehicleDataRow struct {
	VehicleID            string                       `xml:"VEHICLE_ID"`              // 5987
	Route                string                       `xml:"ROUTE"`                   // 1
	RunID                string                       `xml:"RUN_ID"`                  // 21
	TripBlock            string                       `xml:"TRIP_BLOCK"`              // 001HL064
	PatternID            string                       `xml:"PATTERN_ID"`              // 264
	Destination          string                       `xml:"DESTINATION"`             // 1 NEWARK-IVY HILL VIA RIVER TERM
	Longitude            string                       `xml:"LONGITUDE"`               // -74.24513778686523
	Latitude             string                       `xml:"LATITUDE"`                // 40.73779029846192
	GPSTimestmp          string                       `xml:"GPS_TIMESTMP"`            // 25-Apr-2019 12:15:12 AM
	LastModified         string                       `xml:"LAST_MODIFIED"`           // 25-Apr-2019 12:16:10 AM
	AsInternalTripNumber string                       `xml:"AS_INTERNAL_TRIP_NUMBER"` // 13734490
	Timepoints           []BusVehicleDataRowTimepoint `xml:"TIMEPOINTS"`
}

// BusVehicleDataRowTimepoint is part of the BusVehicleDataRow
type BusVehicleDataRowTimepoint struct {
	AsTimingPointID string `xml:"AS_TIMING_POINT_ID"` // IVY HILL
	AsDescription   string `xml:"AS_DESCRIPTION"`     // IVY HILL LOOP (MT VERNON PL)
	AsSchedDepTime  string `xml:"AS_SCHED_DEP_TIME"`  // 25-Apr-2019 12:20:00 AM
}

// ScheduleXGTFSTrip is part of the GetScheduleXGTFSResponse
type ScheduleXGTFSTrip struct {
	GTFSTripID             int    `xml:"gtfs_trip_id"`
	GTFSStopID             int    `xml:"gtfs_stop_id"`
	GTFSStopCode           string `xml:"gtfs_stop_Code"`
	GTFSRouteID            int    `xml:"gtfs_route_id"`
	GTFSServiceID          int    `xml:"gtfs_service_id"`
	GTFSFileDate           string `xml:"gtfs_file_Date"`
	Route                  string `xml:"route"`
	Lanegate               string `xml:"lanegate"` // terminal Line or Gate (".." for inbound buses)
	ScheduledLaneGate      string `xml:"scheduled_lane_gate"`
	ManualLaneGate         string `xml:"manual_lane_gate"`
	DepartureTime          string `xml:"departuretime"`          // example: 1:41 AM
	ScheduledDepartureDate string `xml:"scheduleddeparturedate"` // example: 03-JAN-19
	ScheduledDepartureTime string `xml:"scheduleddeparturetime"` // example: 1:41 AM
	SchedDepTime           string `xml:"sched_dep_Time"`         // example: 03-JAN-19 01.41.00.000000 AM
	SecLate                string `xml:"sec_late"`
	BusHeader              string `xml:"busheader"`
	RunID                  string `xml:"run_id"`
	StopName               string `xml:"stopname"`
	StopCity               string `xml:"stopcity"` // example: NEW YORK CITY
	TripBlock              string `xml:"trip_block"`
	Direction              string `xml:"direction"` // can be "In" or "Ou"
}

// GetNextTripsRequest represents GetNextTrips API request
type GetNextTripsRequest struct {
	StopID string
}

// GetNextTripsResponse represents GetNextTrips API response
type GetNextTripsResponse struct {
	Trips []GetNextTrip `xml:"Trip"`
}

// GetNextTrip represents part of the GetNextTripsResponse
type GetNextTrip struct {
	TripID        string  `xml:"Trip_id"`         // 35971
	ArrivalTime   string  `xml:"arrival_time"`    // 23:00:48
	DepartureTime string  `xml:"departure_time"`  // 23:00:48
	SchedDepTime  string  `xml:"sched_dep_time"`  // 4/22/2019 11:09:00 PM
	StopID        string  `xml:"stop_id"`         // 21884
	StopSequence  string  `xml:"stop_sequence"`   // 84
	Route         string  `xml:"route"`           // 94
	Header        string  `xml:"header"`          // BLOOMFIELD CENTER
	StopName      string  `xml:"stop_name"`       // HESSIAN AVE AT RED BANK AVE#
	TimingPointID string  `xml:"timing_point_id"` // BLFDMUNI
	StopLat       float32 `xml:"stop_lat"`        // 39.862620
	StopLon       float32 `xml:"stop_lon"`        // -75.168910
	SecLate       int     `xml:"sec_late"`        // -60
}

// GetBusDVRequest represents GetBusDV API request
type GetBusDVRequest struct {
	Location string
}

// GetBusDVResponse represents GetBusDV API response
type GetBusDVResponse struct {
	Trip []BusDVTrip `xml:"Trip"`
}

// BusDVTrip represents part of GetBusDVResponse
type BusDVTrip struct {
	PublicRoute   string `xml:"public_route"`  // 123
	Header        string `xml:"header"`        // JERSEY CITY CHRIST HOSP
	Lanegate      string `xml:"lanegate"`      // 303
	DepartureTime string `xml:"departuretime"` // Approaching
	Remarks       string `xml:"remarks"`       //
}

// GetBusLocationsResponse represents GetBusLocations API response
type GetBusLocationsResponse struct {
	Terminal []string `xml:"terminal"`
}

// GetMessagesRequest represents GetMessages API response
type GetMessagesRequest struct {
	StopID string
}

// GetMessagesResponse represents GetMessages API response
type GetMessagesResponse struct {
	Message []string `xml:"message"`
}

// GetScheduleDataRequest represents GetScheduleData API response
type GetScheduleDataRequest struct {
	Site    string
	Minutes string
}

// GetScheduleDataResponse represents GetScheduleData API response
type GetScheduleDataResponse struct {
	Trips []GetScheduleDataTrip `xml:"trop"`
}

// GetScheduleDataTrip represents part of GetScheduleDataResponse
type GetScheduleDataTrip struct {
	InternalTripNumber string                  `xml:"internal_trip_number"` // 13694960
	Route              string                  `xml:"route"`                // 123
	BusHeader          string                  `xml:"BusHeader"`            // Jersey City Christ Hosp
	RunID              string                  `xml:"run_id"`               // 624
	ManualLaneGate     string                  `xml:"manual_lane_gate"`     //
	LaneGate           string                  `xml:"LaneGate"`             // 214-3
	SecLate            string                  `xml:"sec_late"`             //
	Remarks            string                  `xml:"Remarks"`              //
	DepartureTime      string                  `xml:"DepartureTime"`        // 8:14 PM
	Stop               GetScheduleDataTripStop `xml:"STOP"`
}

// GetScheduleDataTripStop represents part of GetScheduleDataTrip
type GetScheduleDataTripStop struct {
	ScheduledDepartureDate string `xml:"scheduleddeparturedate"` // 4/24/2019 12:00:00 AM
	ScheduledDepartureTime string `xml:"cheduleddeparturetime"`  // 8:18 PM
	TopName                string `xml:"topname"`                // Port Authority Bus Terminal
	TopCity                string `xml:"topcity"`                // NEW YORK CITY
}
