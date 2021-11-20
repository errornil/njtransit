package njtransit

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type httpClientMock struct {
	mock.Mock
}

func (h *httpClientMock) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	args := h.Called(url, data)
	return args.Get(0).(*http.Response), args.Error(1)
}

type closingBuffer struct {
	*bytes.Buffer
}

func (cb *closingBuffer) Close() error {
	return nil
}

func TestTrainClientGetStationList(t *testing.T) {
	httpClient := new(httpClientMock)
	httpClient.
		On(
			"PostForm",
			"http://traindata.njtransit.com:8092/NJTTrainData.asmx/getStationListXML",
			url.Values{"username": []string{"username"}, "password": []string{"password"}},
		).
		Return(
			&http.Response{
				StatusCode: http.StatusOK,
				Header: map[string][]string{
					"Content-Type": {"text/xml; charset=utf-8"},
				},
				Body: &closingBuffer{
					bytes.NewBufferString(
						`<?xml version="1.0" encoding="utf-8"?>
						<STATIONS>
							<STATION>
								<STATION_2CHAR>AB</STATION_2CHAR>
								<STATIONNAME>Absecon</STATIONNAME>
							</STATION>
							<STATION>
								<STATION_2CHAR>AZ</STATION_2CHAR>
								<STATIONNAME>Allendale</STATIONNAME>
							</STATION>
						</STATIONS>`,
					),
				},
			},
			nil,
		)

	trainClient := NewTrainDataClient(
		httpClient,
		"username",
		"password",
		"http://traindata.njtransit.com:8092/NJTTrainData.asmx",
	)
	resp, err := trainClient.GetStationList()

	assert.NoError(t, err)
	assert.Equal(
		t,
		&GetStationListResponse{
			Stations: []GetStationListResponseStation{
				{
					TwoChar: "AB",
					Name:    "Absecon",
				},
				{
					TwoChar: "AZ",
					Name:    "Allendale",
				},
			},
		},
		resp,
	)
}

func TestTrainClientGetStationSchedule(t *testing.T) {
	httpClient := new(httpClientMock)
	httpClient.
		On(
			"PostForm",
			"http://traindata.njtransit.com:8092/NJTTrainData.asmx/getStationScheduleXML",
			url.Values{
				"username": []string{"username"},
				"password": []string{"password"},
				"station":  []string{"AZ"},
				"NJT_Only": []string{"1"},
			},
		).
		Return(
			&http.Response{
				StatusCode: http.StatusOK,
				Header: map[string][]string{
					"Content-Type": {"text/xml; charset=utf-8"},
				},
				Body: &closingBuffer{
					bytes.NewBufferString(
						`<STATION>
							<STATION_2CHAR>AZ</STATION_2CHAR>
							<STATIONNAME>Allendale</STATIONNAME>
							<ITEMS>
								<ITEM>
									<ITEM_INDEX>0</ITEM_INDEX>
									<SCHED_DEP_DATE>10-Sep-2019 12:12:15 AM</SCHED_DEP_DATE>
									<DESTINATION>Suffern</DESTINATION>
									<SCHED_TRACK>1</SCHED_TRACK>
									<TRAIN_ID>1183</TRAIN_ID>
									<LINE>Bergen County Line</LINE>
									<STATION_POSITION>1</STATION_POSITION>
									<DIRECTION>Westbound</DIRECTION>
									<DWELL_TIME>45</DWELL_TIME>
									<PERM_CONNECTING_TRAIN_ID></PERM_CONNECTING_TRAIN_ID>
									<PERM_PICKUP>False</PERM_PICKUP>
									<PERM_DROPOFF></PERM_DROPOFF>
									<STOP_CODE>S</STOP_CODE>
									<STOPPING_AT>Ramsey,Ramsey Rt 17,Mahwah,Suffern</STOPPING_AT>
								</ITEM>
							</ITEMS>
						</STATION>`,
					),
				},
			},
			nil,
		)

	trainClient := NewTrainDataClient(
		httpClient,
		"username",
		"password",
		"http://traindata.njtransit.com:8092/NJTTrainData.asmx",
	)
	resp, err := trainClient.GetStationSchedule("AZ", true)

	assert.NoError(t, err)
	assert.Equal(
		t,
		&GetStationScheduleResponse{
			TwoChar: "AZ",
			Name:    "Allendale",
			Items: []GetStationScheduleResponseItem{
				{
					ItemIndex:             0,
					SchedDepDate:          "10-Sep-2019 12:12:15 AM",
					Destination:           "Suffern",
					SchedTrack:            "1",
					TrainID:               "1183",
					Line:                  "Bergen County Line",
					StationPosition:       "1",
					Direction:             "Westbound",
					DwellTimeSeconds:      45,
					PermConnectingTrainID: "",
					PermPickup:            "False",
					PermDropoff:           "",
					StopCode:              "S",
					StoppingAt:            "Ramsey,Ramsey Rt 17,Mahwah,Suffern",
				},
			},
		},
		resp,
	)
}

func TestTrainClientGetStationMessage(t *testing.T) {
	httpClient := new(httpClientMock)
	httpClient.
		On(
			"PostForm",
			"http://traindata.njtransit.com:8092/NJTTrainData.asmx/getStationMSGXML",
			url.Values{
				"username":  []string{"username"},
				"password":  []string{"password"},
				"station":   []string{"NY"},
				"trainLine": []string{""},
			},
		).
		Return(
			&http.Response{
				StatusCode: http.StatusOK,
				Header: map[string][]string{
					"Content-Type": {"text/xml; charset=utf-8"},
				},
				Body: &closingBuffer{
					bytes.NewBufferString(
						`<STATION>
							<STATION_2CHAR>NY</STATION_2CHAR>
							<STATIONNAME>New York</STATIONNAME>
							<BANNERMSGS/>
							<ITEMS>
								<ITEM>
									<ITEM_INDEX>0</ITEM_INDEX>
									<SCHED_DEP_DATE>11-Sep-2019 12:02:00 AM</SCHED_DEP_DATE>
									<DESTINATION>Dover -SEC</DESTINATION>
									<TRACK>4</TRACK>
									<LINE>Morristown Line</LINE>
									<TRAIN_ID>6683</TRAIN_ID>
									<CONNECTING_TRAIN_ID></CONNECTING_TRAIN_ID>
									<STATUS>ALL ABOARD</STATUS>
									<SEC_LATE>-60</SEC_LATE>
									<LAST_MODIFIED>10-Sep-2019 11:51:58 PM</LAST_MODIFIED>
									<BACKCOLOR>green</BACKCOLOR>
									<FORECOLOR>white</FORECOLOR>
									<SHADOWCOLOR>black</SHADOWCOLOR>
									<GPSLATITUDE></GPSLATITUDE>
									<GPSLONGITUDE></GPSLONGITUDE>
									<GPSTIME>11-Sep-2019 12:00:02 AM</GPSTIME>
									<STATION_POSITION>0</STATION_POSITION>
									<LINEABBREVIATION>M&E</LINEABBREVIATION>
									<INLINEMSG></INLINEMSG>
								</ITEM>
							</ITEMS>
						</STATION>`,
					),
				},
			},
			nil,
		)

	trainClient := NewTrainDataClient(
		httpClient,
		"username",
		"password",
		"http://traindata.njtransit.com:8092/NJTTrainData.asmx",
	)
	resp, err := trainClient.GetStationMessage("NY", "")

	assert.NoError(t, err)
	assert.Equal(
		t,
		&GetStationMessageResponse{
			TwoChar: "NY",
			Name:    "New York",
			Items: []*GetStationMessageResponseItem{
				{
					ItemIndex:         "0",
					SchedDepDate:      "11-Sep-2019 12:02:00 AM",
					Destination:       "Dover",
					Track:             "4",
					Line:              "Morristown Line",
					TrainID:           "6683",
					ConnectingTrainID: "",
					Status:            "ALL ABOARD",
					SecLate:           "-60",
					LastModified:      "10-Sep-2019 11:51:58 PM",
					BackgroundColor:   "green",
					ForegroundColor:   "white",
					ShadowColor:       "black",
					GPSLatitude:       "",
					GPSLongitude:      "",
					GPSTime:           "11-Sep-2019 12:00:02 AM",
					StationPosition:   "0",
					LineAbbreviation:  "M&E",
					InlineMessage:     "",
				},
			},
		},
		resp,
	)
}

func TestTrainClientGetTrainSchedule19Rec(t *testing.T) {
	httpClient := new(httpClientMock)
	httpClient.
		On(
			"PostForm",
			"http://traindata.njtransit.com:8092/NJTTrainData.asmx/getTrainScheduleXML19Rec",
			url.Values{
				"username": []string{"username"},
				"password": []string{"password"},
				"station":  []string{"NY"},
			},
		).
		Return(
			&http.Response{
				StatusCode: http.StatusOK,
				Header: map[string][]string{
					"Content-Type": {"text/xml; charset=utf-8"},
				},
				Body: &closingBuffer{
					bytes.NewBufferString(
						`<STATION>
							<STATION_2CHAR>RH</STATION_2CHAR>
							<STATIONNAME>Rahway</STATIONNAME>
							<BANNERMSGS/>
							<ITEMS>
							<ITEM>
								<ITEM_INDEX>0</ITEM_INDEX>
								<SCHED_DEP_DATE>12-Oct-2019 11:46:30 PM</SCHED_DEP_DATE>
								<DESTINATION>Long Branch-BH</DESTINATION>
								<TRACK>B</TRACK>
								<LINE>No Jersey Coast</LINE>
								<TRAIN_ID>7285</TRAIN_ID>
								<CONNECTING_TRAIN_ID>4785</CONNECTING_TRAIN_ID>
								<STATUS>in 24 Min</STATUS>
								<SEC_LATE>534</SEC_LATE>
								<LAST_MODIFIED>12-Oct-2019 11:29:46 PM</LAST_MODIFIED>
								<BACKCOLOR>CornflowerBlue</BACKCOLOR>
								<FORECOLOR>white</FORECOLOR>
								<SHADOWCOLOR>black</SHADOWCOLOR>
								<GPSLATITUDE>40.7354</GPSLATITUDE>
								<GPSLONGITUDE>-74.1632</GPSLONGITUDE>
								<GPSTIME>12-Oct-2019 11:29:45 PM</GPSTIME>
								<STATION_POSITION>1</STATION_POSITION>
								<LINEABBREVIATION>NJCL</LINEABBREVIATION>
								<INLINEMSG></INLINEMSG>
							</ITEM>
						</STATION>`,
					),
				},
			},
			nil,
		)

	trainClient := NewTrainDataClient(
		httpClient,
		"username",
		"password",
		"http://traindata.njtransit.com:8092/NJTTrainData.asmx",
	)
	resp, err := trainClient.GetTrainSchedule19Rec("NY")

	assert.NoError(t, err)
	assert.Equal(
		t,
		&GetTrainSchedule19RecResponse{
			TwoChar: "RH",
			Name:    "Rahway",
			Items: []*GetTrainSchedule19RecResponseItem{
				{
					ItemIndex:         0,
					SchedDepDate:      "12-Oct-2019 11:46:30 PM",
					Destination:       "Long Branch",
					Track:             "B",
					Line:              "No Jersey Coast",
					TrainID:           "7285",
					ConnectingTrainID: "4785",
					Status:            "in 24 Min",
					SecLate:           534,
					LastModified:      "12-Oct-2019 11:29:46 PM",
					BackgroundColor:   "CornflowerBlue",
					ForegroundColor:   "white",
					ShadowColor:       "black",
					GPSLatitude:       "40.7354",
					GPSLongitude:      "-74.1632",
					GPSTime:           "12-Oct-2019 11:29:45 PM",
					StationPosition:   "1",
					LineAbbreviation:  "NJCL",
					InlineMessage:     "",
				},
			},
		},
		resp,
	)
}

func TestTrainClientGetTrainSchedule19RecInvalidResponse(t *testing.T) {
	httpClient := new(httpClientMock)
	httpClient.
		On(
			"PostForm",
			"http://traindata.njtransit.com:8092/NJTTrainData.asmx/getTrainScheduleXML19Rec",
			url.Values{
				"username": []string{"username"},
				"password": []string{"password"},
				"station":  []string{"NY"},
			},
		).
		Return(
			&http.Response{
				StatusCode: http.StatusOK,
				Header: map[string][]string{
					"Content-Type": {"text/html; charset=utf-8"},
				},
				Body: &closingBuffer{
					bytes.NewBufferString(
						`

						<!DOCTYPE html>
						Please enter credentials. <a href="NJTTrainData.asmx">Click here</a> to go back.`,
					),
				},
			},
			nil,
		)

	trainClient := NewTrainDataClient(
		httpClient,
		"username",
		"password",
		"http://traindata.njtransit.com:8092/NJTTrainData.asmx",
	)
	_, err := trainClient.GetTrainSchedule19Rec("NY")

	assert.Error(t, err)
}
