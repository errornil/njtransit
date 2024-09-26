package njtransit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	gtfs "github.com/errornil/njtransit/proto/transit_realtime"
	"google.golang.org/protobuf/proto"
)

const (
	BusProdURL = "https://pcsdata.njtransit.com/api/GTFS/"
	BusTestURL = "https://testpcsdata.njtransit.com/api/GTFS/"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// BusClient holds information between API calls
type BusClient struct {
	url       string
	username  string
	password  string
	userAgent string
	client    HTTPClient

	token string
}

// NewBusClient creates new BusClient
func NewBusClient(
	url,
	username,
	password,
	userAgent string,
	client HTTPClient,
) (*BusClient, error) {
	bc := &BusClient{
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
func (bc *BusClient) AuthenticateUser() error {
	// set username and password as data-url-encoded
	body := url.Values{}
	body.Add("username", bc.username)
	body.Add("password", bc.password)
	b := body.Encode()

	url := bc.url + "authenticateUser"

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

func (bc *BusClient) GetGTFS() ([]byte, error) {
	b, err := bc.callAPI("getGTFS")
	if err != nil {
		return nil, fmt.Errorf("callAPI: %v", err)
	}

	return b, nil
}

func (bc *BusClient) GetTripUpdates() (*gtfs.FeedMessage, error) {
	return bc.callAPIProto("getTripUpdates")
}

func (bc *BusClient) GetVehiclePositions() (*gtfs.FeedMessage, error) {
	return bc.callAPIProto("getVehiclePositions")
}

func (bc *BusClient) GetAlerts() (*gtfs.FeedMessage, error) {
	return bc.callAPIProto("getAlerts")
}

func (bc *BusClient) callAPI(url string) ([]byte, error) {
	reqBody := &bytes.Buffer{}
	writer := multipart.NewWriter(reqBody)
	err := writer.WriteField("token", bc.token)
	if err != nil {
		return nil, fmt.Errorf("write: %v", err)
	}
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("close writer: %v", err)
	}

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

func (bc *BusClient) callAPIProto(url string) (*gtfs.FeedMessage, error) {
	b, err := bc.callAPI(url)
	if err != nil {
		return nil, fmt.Errorf("callAPI: %v", err)
	}

	feed := &gtfs.FeedMessage{}
	err = proto.Unmarshal(b, feed)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response: %v", err)
	}

	return feed, nil
}
