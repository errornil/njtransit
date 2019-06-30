package njtransit

import (
	"encoding/xml"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Web Serivces URLs
const (
	BusDataProdURL = "https://busdata.njtransit.com/NJTBusData.asmx"
	BusDataTestURL = "https://busdata_tst.njtransit.com/NJTBusData.asmx"
)

// Pre-defined locations
// That can be used in GetBusDV method
const (
	CamdenBusTerminal        = "CAMD"
	AtlanticCityBusTerminal  = "ATLC"
	HackensackBusTerminal    = "HACK"
	HobokenBusTerminal       = "HBKN"
	IrvingtonBusTerminal     = "IRVN"
	LakewoodBusTerminal      = "LKWD"
	NewarkPennStation        = "NWRK"
	MetroParkBusTerminal     = "MTPK"
	PortAuthorityBusTerminal = "PABT"
	OldBridgeBusTerminal     = "OBRG"
	WayneBusTerminal         = "WYNE"
	TrentonBusTerminal       = "TREN"
	GeorgeWashingtonTerminal = "GWBT"
)

// BusDataClient holds information between API calls
type BusDataClient struct {
	username   string
	password   string
	busDataURL string

	// Map from VehicleID to BusVehicleData checksum
	// used in GetBusVehicleDataStream to dedupe messages
	busVehicleDataRowsChecksum map[string]uint32
}

// NewBusDataClient creates new BusDataClient
func NewBusDataClient(username, password, busDataURL string) *BusDataClient {
	return &BusDataClient{
		username:   username,
		password:   password,
		busDataURL: busDataURL,
	}
}

// GetBusVehicleData - Status By Bus data
// This Method will provide Bus Vehicle Information.
// It will list the vehicles currently reporting real-time information.
func (c *BusDataClient) GetBusVehicleData() (*GetBusVehicleDataResponse, error) {
	v := url.Values{}
	v.Add("username", c.username)
	v.Add("password", c.password)

	resp, err := http.PostForm(fmt.Sprintf("%s/getBusVehicleDataXML", c.busDataURL), v)
	if err != nil {
		return nil, fmt.Errorf("failed to send GetBusVehicleData request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read GetBusVehicleData response: %v", err)
	}

	response := GetBusVehicleDataResponse{}
	err = xml.Unmarshal(body, &response)

	if err != nil {
		return nil, fmt.Errorf("failed to parse GetBusVehicleData response: %v, body: %s", err, body)
	}

	return &response, nil
}

// GetBusVehicleDataStream provides a stream or BusVehicleData updates.
func (c *BusDataClient) GetBusVehicleDataStream(r chan BusVehicleDataRow, e chan error, updateTnterval time.Duration, dedupe bool) {
	c.busVehicleDataRowsChecksum = map[string]uint32{}

	for {
		resp, err := c.GetBusVehicleData()
		if err != nil {
			e <- err
			continue
		}

		for _, row := range resp.Rows {
			if !dedupe || c.isUniqueBusVehicleDataRow(row) {
				r <- row
			}
		}
		time.Sleep(updateTnterval)
	}
}

// GetNextTrips retrieves upcoming bus arrivals for given stopID
// This method provides schedule information.
// The data consist of the next 20 trips that depart the given stop
func (c *BusDataClient) GetNextTrips(request GetNextTripsRequest) (*GetNextTripsResponse, error) {
	v := url.Values{}
	v.Add("username", c.username)
	v.Add("password", c.password)
	v.Add("stopid", fmt.Sprintf("%d", request.StopID))

	resp, err := http.PostForm(fmt.Sprintf("%s/getNextTripsXML", c.busDataURL), v)
	if err != nil {
		return nil, fmt.Errorf("failed to send GetNextTrips request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read GetNextTrips response: %v", err)
	}

	response := GetNextTripsResponse{}
	err = xml.Unmarshal(body, &response)

	if err != nil {
		return nil, fmt.Errorf("failed to parse GetNextTrips response: %v, body: %s", err, body)
	}

	return &response, nil
}

// GetBusDV - Gets the first layer of data for BusDV
// This method provides schedule information.
// The data consist of the next trip to depart for each possible lane at the requested location.
func (c *BusDataClient) GetBusDV(request GetBusDVRequest) (*GetBusDVResponse, error) {
	v := url.Values{}
	v.Add("username", c.username)
	v.Add("password", c.password)
	v.Add("location", request.Location)

	resp, err := http.PostForm(fmt.Sprintf("%s/getBusDVXML", c.busDataURL), v)
	if err != nil {
		return nil, fmt.Errorf("failed to send GetBusDV request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read GetBusDV response: %v", err)
	}

	response := GetBusDVResponse{}
	err = xml.Unmarshal(body, &response)

	if err != nil {
		return nil, fmt.Errorf("failed to parse GetBusDV response: %v, body: %s", err, body)
	}

	return &response, nil
}

// GetBusLocations - This method provides a list of locations that can be used in the GetBusDVXML
func (c *BusDataClient) GetBusLocations() (*GetBusLocationsResponse, error) {
	v := url.Values{}
	v.Add("username", c.username)
	v.Add("password", c.password)

	resp, err := http.PostForm(fmt.Sprintf("%s/getBusLocationsXML", c.busDataURL), v)
	if err != nil {
		return nil, fmt.Errorf("failed to send GetBusLocations request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read GetBusLocations response: %v", err)
	}

	response := GetBusLocationsResponse{}
	err = xml.Unmarshal(body, &response)

	if err != nil {
		return nil, fmt.Errorf("failed to parse GetBusLocations response: %v, body: %s", err, body)
	}

	return &response, nil
}

// GetMessages - This method provides a list of messages
func (c *BusDataClient) GetMessages(request GetMessagesRequest) (*GetMessagesResponse, error) {
	v := url.Values{}
	v.Add("username", c.username)
	v.Add("password", c.password)
	v.Add("stopid", fmt.Sprintf("%d", request.StopID))

	resp, err := http.PostForm(fmt.Sprintf("%s/getMessagesXML", c.busDataURL), v)
	if err != nil {
		return nil, fmt.Errorf("failed to send GetMessages request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read GetMessages response: %v", err)
	}

	response := GetMessagesResponse{}
	err = xml.Unmarshal(body, &response)

	if err != nil {
		return nil, fmt.Errorf("failed to parse GetMessages response: %v, body: %s", err, body)
	}

	return &response, nil
}

// GetScheduleData - This method will provide schedule information.
// The data consists of the departures for the given site that depart within the given number of minutes.
// Also included are the remaining stops that each of these trips will be making.
func (c *BusDataClient) GetScheduleData(request GetScheduleDataRequest) (*GetScheduleDataResponse, error) {
	v := url.Values{}
	v.Add("username", c.username)
	v.Add("password", c.password)
	v.Add("site", request.Site)
	v.Add("minutes", fmt.Sprintf("%d", request.Minutes))

	resp, err := http.PostForm(fmt.Sprintf("%s/getScheduleDataXML", c.busDataURL), v)
	if err != nil {
		return nil, fmt.Errorf("failed to send GetScheduleData request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read GetScheduleData response: %v", err)
	}

	response := GetScheduleDataResponse{}
	err = xml.Unmarshal(body, &response)

	if err != nil {
		return nil, fmt.Errorf("failed to parse GetScheduleData response: %v, body: %s", err, body)
	}

	return &response, nil
}

// GetScheduleXGTFS -This method will provide schedule information.
// The data consists of the arrivals and departures for the given site that depart within the given number of minutes.
func (c *BusDataClient) GetScheduleXGTFS(request GetScheduleXGTFSRequest) (*GetScheduleXGTFSResponse, error) {
	v := url.Values{}
	v.Add("username", c.username)
	v.Add("password", c.password)
	v.Add("site", request.Site)
	v.Add("minutes", fmt.Sprintf("%d", request.Minutes))

	resp, err := http.PostForm(fmt.Sprintf("%s/getScheduleXGTFS", c.busDataURL), v)
	if err != nil {
		return nil, fmt.Errorf("failed to send GetScheduleXGTFS request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read GetScheduleXGTFS response: %v", err)
	}

	response := GetScheduleXGTFSResponse{}
	err = xml.Unmarshal(body, &response)

	if err != nil {
		return nil, fmt.Errorf("failed to parse GetScheduleXGTFS response: %v, body: %s", err, body)
	}

	return &response, nil
}

func (c *BusDataClient) isUniqueBusVehicleDataRow(row BusVehicleDataRow) bool {
	newHash := busVehicleDataRowChecksum(row)
	currentHash, ok := c.busVehicleDataRowsChecksum[row.VehicleID]

	if !ok {
		c.busVehicleDataRowsChecksum[row.VehicleID] = newHash
		return true
	}

	if currentHash != newHash {
		c.busVehicleDataRowsChecksum[row.VehicleID] = newHash
		return true
	}

	return false
}

func busVehicleDataRowChecksum(row BusVehicleDataRow) uint32 {
	return crc32.ChecksumIEEE([]byte(fmt.Sprintf("%#v", row)))
}
