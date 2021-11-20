package njtransit

type GetStationListResponse struct {
	Stations []GetStationListResponseStation `xml:"STATION"`
}

type GetStationListResponseStation struct {
	TwoChar string `xml:"STATION_2CHAR"`
	Name    string `xml:"STATIONNAME"`
}

type GetStationScheduleResponse struct {
	TwoChar string                           `xml:"STATION_2CHAR"`
	Name    string                           `xml:"STATIONNAME"`
	Items   []GetStationScheduleResponseItem `xml:"ITEMS>ITEM"`
}

type GetStationScheduleResponseItem struct {
	ItemIndex             int    `xml:"ITEM_INDEX"`
	SchedDepDate          string `xml:"SCHED_DEP_DATE"`
	Destination           string `xml:"DESTINATION"`
	SchedTrack            string `xml:"SCHED_TRACK"` // see Appendix II
	TrainID               string `xml:"TRAIN_ID"`
	Line                  string `xml:"LINE"`
	StationPosition       string `xml:"STATION_POSITION"` // see Appendix III
	Direction             string `xml:"DIRECTION"`
	DwellTimeSeconds      int    `xml:"DWELL_TIME"`
	PermConnectingTrainID string `xml:"PERM_CONNECTING_TRAIN_ID"`
	PermPickup            string `xml:"PERM_PICKUP"`
	PermDropoff           string `xml:"PERM_DROPOFF"`
	StopCode              string `xml:"STOP_CODE"` // see Appendix IV
	StoppingAt            string `xml:"STOPPING_AT"`
}

type GetStationMessageResponse struct {
	TwoChar       string                           `xml:"STATION_2CHAR"`
	Name          string                           `xml:"STATIONNAME"`
	BannerMessage string                           `xml:"BANNERMSGS"`
	Items         []*GetStationMessageResponseItem `xml:"ITEMS>ITEM"`
}

type GetStationMessageResponseItem struct {
	ItemIndex         string `xml:"ITEM_INDEX"`
	SchedDepDate      string `xml:"SCHED_DEP_DATE"`
	Destination       string `xml:"DESTINATION"`
	Track             string `xml:"TRACK"`
	Line              string `xml:"LINE"`
	TrainID           string `xml:"TRAIN_ID"`
	ConnectingTrainID string `xml:"CONNECTING_TRAIN_ID"`
	Status            string `xml:"STATUS"`
	SecLate           string `xml:"SEC_LATE"`
	LastModified      string `xml:"LAST_MODIFIED"`
	BackgroundColor   string `xml:"BACKCOLOR"`
	ForegroundColor   string `xml:"FORECOLOR"`
	ShadowColor       string `xml:"SHADOWCOLOR"`
	GPSLatitude       string `xml:"GPSLATITUDE"`
	GPSLongitude      string `xml:"GPSLONGITUDE"`
	GPSTime           string `xml:"GPSTIME"`
	StationPosition   string `xml:"STATION_POSITION"`
	LineAbbreviation  string `xml:"LINEABBREVIATION"`
	InlineMessage     string `xml:"INLINEMSG"`
}

type GetTrainScheduleResponse struct {
}

type GetTrainSchedule19RecResponse struct {
	TwoChar string                               `xml:"STATION_2CHAR"`
	Name    string                               `xml:"STATIONNAME"`
	Items   []*GetTrainSchedule19RecResponseItem `xml:"ITEMS>ITEM"`
}

type GetTrainSchedule19RecResponseItem struct {
	ItemIndex         int    `xml:"ITEM_INDEX"`
	SchedDepDate      string `xml:"SCHED_DEP_DATE"`
	Destination       string `xml:"DESTINATION"`
	Track             string `xml:"TRACK"`
	Line              string `xml:"LINE"`
	TrainID           string `xml:"TRAIN_ID"`
	ConnectingTrainID string `xml:"CONNECTING_TRAIN_ID"`
	Status            string `xml:"STATUS"`
	SecLate           int    `xml:"SEC_LATE"`
	LastModified      string `xml:"LAST_MODIFIED"`
	BackgroundColor   string `xml:"BACKCOLOR"`
	ForegroundColor   string `xml:"FORECOLOR"`
	ShadowColor       string `xml:"SHADOWCOLOR"`
	GPSLatitude       string `xml:"GPSLATITUDE"`
	GPSLongitude      string `xml:"GPSLONGITUDE"`
	GPSTime           string `xml:"GPSTIME"`
	StationPosition   string `xml:"STATION_POSITION"`
	LineAbbreviation  string `xml:"LINEABBREVIATION"`
	InlineMessage     string `xml:"INLINEMSG"`
}

type GetVehicleDataResponse struct {
}
