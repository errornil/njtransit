package njtransit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

const (
	BusDVProdURL = "https://pcsdata.njtransit.com/api/BUSDV2/"
	BusDVTestURL = "https://testpcsdata.njtransit.com/api/BUSDV2/"
)

// BusDV2Client holds information between API calls
type BusDV2Client struct {
	url       string
	username  string
	password  string
	userAgent string
	client    HTTPClient

	token string
}

type DVTrip struct {
	PublicRoute     string `json:"public_route"`
	Header          string `json:"header"`
	Lanegate        string `json:"lanegate"`
	DepartureTime   string `json:"departuretime"`
	Remarks         string `json:"remarks"`
	InternalTripNum string `json:"internal_trip_number"`
	SchedDepTime    string `json:"sched_dep_time"`
	TimingPointID   string `json:"timing_point_id"`
	Message         string `json:"message"`
	FullScreen      string `json:"fullscreen"`
	PassLoad        string `json:"passload"`
	VehicleID       string `json:"vehicle_id"`
}

type VehicleLocation struct {
	VehicleLat                string `json:"VehicleLat"`
	VehicleLong               string `json:"VehicleLong"`
	VehicleID                 string `json:"VehicleID"`
	VehiclePassengerLoad      string `json:"VehiclePassengerLoad"`
	VehicleRoute              string `json:"VehicleRoute"`
	VehicleDestination        string `json:"VehicleDestination"`
	VehicleDistanceMiles      string `json:"VehicleDistanceMiles"`
	VehicleInternalTripNumber string `json:"VehicleInternalTripNumber"`
	VehicleScheduledDeparture string `json:"VehicleScheduledDeparture"`
}

type GetBusDVResponse struct {
	Message struct {
		Message string `json:"message"`
	} `json:"message"`
	DVTrip DVTrip `json:"DVTrip"`
}

type GetVehicleLocations []VehicleLocation

// NewBusDV2Client creates new BusDV2Client
func NewBusDV2Client(
	url,
	username,
	password,
	userAgent string,
	client HTTPClient,
) (*BusDV2Client, error) {
	bc := &BusDV2Client{
		url:       url,
		username:  username,
		password:  password,
		userAgent: userAgent,
		client:    client,
	}

	err := bc.AuthenticateUser()
	if err != nil {
		return nil, err
	}

	return bc, nil
}

// AuthenticateUser gets token for user
func (bc *BusDV2Client) AuthenticateUser() error {
	// set username and password as data-url-encoded
	body := url.Values{}
	body.Add("username", bc.username)
	body.Add("password", bc.password)
	b := body.Encode()

	url := bc.url + "authenticateUser"
	log.Printf("→ %s", url)

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(b))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Set("User-Agent", bc.userAgent)

	resp, err := bc.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	respb := bytes.Buffer{}
	_, err = io.Copy(&respb, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to authenticate user, status code: %d", resp.StatusCode)
	}

	var response struct {
		Authenticated string `json:"Authenticated"`
		UserToken     string `json:"UserToken"`
	}

	err = json.NewDecoder(&respb).Decode(&response)
	if err != nil {
		return fmt.Errorf("failed to decode response: %v, body: %s", err, respb.String())
	}

	if response.Authenticated != "True" {
		return fmt.Errorf("failed to authenticate user")
	}

	bc.token = response.UserToken
	return nil
}

func (bc *BusDV2Client) GetBusDV(stop, direction, route, ip string) (*GetBusDVResponse, error) {
	var pairs []string
	if stop != "" {
		pairs = append(pairs, "stop", stop)
	}
	if direction != "" {
		pairs = append(pairs, "direction", direction)
	}
	if route != "" {
		pairs = append(pairs, "route", route)
	}
	if ip != "" {
		pairs = append(pairs, "ip", ip)
	}

	response := &GetBusDVResponse{}
	err := bc.callAPIJSON("getBusDV", pairs, response)
	if err != nil {
		return nil, fmt.Errorf("callAPI: %v", err)
	}

	return response, nil
}

func (bc *BusDV2Client) GetVehicleLocations(lat, lon string, radius int, mode string) (*GetVehicleLocations, error) {
	var pairs []string
	if lat != "" {
		pairs = append(pairs, "lat", lat)
	}
	if lon != "" {
		pairs = append(pairs, "lon", lon)
	}
	if radius != 0 {
		pairs = append(pairs, "radius", fmt.Sprintf("%d", radius))
	}
	if mode != "" {
		pairs = append(pairs, "mode", mode)
	}

	response := &GetVehicleLocations{}
	err := bc.callAPIJSON("getVehicleLocations", pairs, response)
	if err != nil {
		return nil, fmt.Errorf("callAPI: %v", err)
	}

	return response, nil
}

func (bc *BusDV2Client) callAPI(url string, bodyPairs []string) ([]byte, error) {
	reqBody := &bytes.Buffer{}
	writer := multipart.NewWriter(reqBody)
	err := writer.WriteField("token", bc.token)
	if err != nil {
		return nil, fmt.Errorf("write: %v", err)
	}

	if len(bodyPairs)%2 != 0 {
		return nil, fmt.Errorf("bodyPairs must be even")
	}
	for i := 0; i < len(bodyPairs); i += 2 {
		err = writer.WriteField(bodyPairs[i], bodyPairs[i+1])
		if err != nil {
			return nil, fmt.Errorf("write: %v", err)
		}
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("close writer: %v", err)
	}

	log.Printf("→ %s", bc.url+url)
	req, err := http.NewRequest(http.MethodPost, bc.url+url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %v", err)
	}

	req.Header.Set("Content-Type", "multipart/form-data; boundary="+writer.Boundary())
	req.Header.Set("User-Agent", bc.userAgent)
	req.Header.Set("Accept", "*/*")

	resp, err := bc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call API: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	body := bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %v", err)
	}

	return body.Bytes(), nil
}

func (bc *BusDV2Client) callAPIJSON(url string, bodyPairs []string, v interface{}) error {
	b, err := bc.callAPI(url, bodyPairs)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, v)
	if err != nil {
		return fmt.Errorf("unmarshal response: %v", err)
	}

	return nil
}
