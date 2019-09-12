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
				Body: &closingBuffer{
					bytes.NewBufferString(
						`<STATIONS>
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
									<PERM_DROPOFF>False</PERM_DROPOFF>
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
					PermPickup:            false,
					PermDropoff:           false,
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
			Items: []GetStationMessageResponseItem{
				{
					ItemIndex:         "0",
					SchedDepDate:      "11-Sep-2019 12:02:00 AM",
					Destination:       "Dover -SEC",
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
