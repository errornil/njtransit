package njtransit

import (
	"encoding/xml"
	"fmt"
	"html"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

const (
	TrainDataProdURL = "https://traindata.njtransit.com/NJTTrainData.asmx"
	TrainDataTestURL = "https://njttraindatatst.njtransit.com/NJTTrainData.asmx"
)

var ErrNotImplemented = errors.New("not implemented")

type TrainDataClient struct {
	httpClient   httpClient
	username     string
	password     string
	trainDataURL string
	replacer     *strings.Replacer
}

func NewTrainDataClient(httpClient httpClient, username, password, trainDataURL string) *TrainDataClient {
	return &TrainDataClient{
		httpClient:   httpClient,
		username:     username,
		password:     password,
		trainDataURL: trainDataURL,
		replacer: strings.NewReplacer(
			" -SEC", "",
			"-BH", "",
		),
	}
}

// GetStationList - List all stations
func (t *TrainDataClient) GetStationList() (*GetStationListResponse, error) {
	v := url.Values{}
	v.Add("username", t.username)
	v.Add("password", t.password)

	resp, err := t.httpClient.PostForm(fmt.Sprintf("%s/getStationListXML", t.trainDataURL), v)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send GetStationList request")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read GetStationList response")
	}

	response := &GetStationListResponse{}
	err = xml.Unmarshal(body, response)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse GetStationList response, body: %s", string(body))
	}

	return response, nil
}

// GetStationSchedule - Provides a list of the 27 hours of train schedule data for any one station or all stations.
// Limited access to 10 times per day but only needed once per day after midnight - 12:30 would be better -
// to show the schedule for the 27 hour period from 12 midnight until 3am the next day.
// The GTFS data does not always match the daily schedules in our train control system.
// NJT_Only is a filter, pass value 1 for NJT trains only; pass value 0 for All trains
func (t *TrainDataClient) GetStationSchedule(station string, njtransitOnly bool) (*GetStationScheduleResponse, error) {
	njtransitOnlyValue := "0"
	if njtransitOnly {
		njtransitOnlyValue = "1"
	}

	v := url.Values{}
	v.Add("username", t.username)
	v.Add("password", t.password)
	v.Add("station", station)
	v.Add("NJT_Only", njtransitOnlyValue)

	resp, err := t.httpClient.PostForm(fmt.Sprintf("%s/getStationScheduleXML", t.trainDataURL), v)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send GetStationSchedule request")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read GetStationSchedule response")
	}

	response := &GetStationScheduleResponse{}
	err = xml.Unmarshal(body, response)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse GetStationSchedule response, body: %s", string(body))
	}

	return response, nil
}

// GetStationMessage - Gets the all station message, but when pass station code,
// returns station message. Note – this is provided by a third party from our above APIs.
func (t *TrainDataClient) GetStationMessage(station, trainLine string) (*GetStationMessageResponse, error) {
	v := url.Values{}
	v.Add("username", t.username)
	v.Add("password", t.password)
	v.Add("station", station)
	v.Add("trainLine", trainLine)

	resp, err := t.httpClient.PostForm(fmt.Sprintf("%s/getStationMSGXML", t.trainDataURL), v)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send GetStationMessage request")
	}
	defer resp.Body.Close()

	decoder := xml.NewDecoder(resp.Body)
	decoder.Strict = false

	response := &GetStationMessageResponse{}
	err = decoder.Decode(response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode GetStationMessage response")
	}

	for _, item := range response.Items {
		item.Destination = strings.TrimSpace(
			t.replacer.Replace(
				html.UnescapeString(item.Destination),
			),
		)
	}

	return response, nil
}

// GetTrainSchedule - List train schedule for a given station,
// data is much the same as DepartureVision with train stop list information
func (t *TrainDataClient) GetTrainSchedule(station string, njtransitOnly bool) (*GetTrainScheduleResponse, error) {
	return nil, ErrNotImplemented
}

// GetTrainScheduleJSON19Rec - List train schedule for a given station,
// data is much the same as DepartureVision, but without train stop list information.
func (t *TrainDataClient) GetTrainSchedule19Rec(station string) (*GetTrainSchedule19RecResponse, error) {
	v := url.Values{}
	v.Add("username", t.username)
	v.Add("password", t.password)
	v.Add("station", station)

	resp, err := t.httpClient.PostForm(fmt.Sprintf("%s/getTrainScheduleXML19Rec", t.trainDataURL), v)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send GetTrainSchedule19Rec request")
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if contentType != "text/xml; charset=utf-8" {
		return nil, fmt.Errorf("invalid response Content-Type: %s", contentType)
	}

	decoder := xml.NewDecoder(resp.Body)
	decoder.Strict = false

	response := &GetTrainSchedule19RecResponse{}
	err = decoder.Decode(response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode GetTrainSchedule19Rec response")
	}

	for _, item := range response.Items {
		item.Destination = strings.TrimSpace(
			t.replacer.Replace(
				html.UnescapeString(item.Destination),
			),
		)
	}

	return response, nil
}

// GetVehicleDataXML - Provides the real-time position data for each active train.
// Provides the latest position, next station and seconds late for any train that
// has moved in the last 5 minutes.
// There is a limit of 40,000 requests per day.
func (t *TrainDataClient) GetVehicleData() (*GetVehicleDataResponse, error) {
	return nil, ErrNotImplemented
}

// func (t *TrainDataClient) GetGTFSRealTimeFeed() {}
