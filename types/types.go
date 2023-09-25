package types

type Message struct {
	Event     string
	ErrorData string
	DateTime  int64
	//
	IMEI string
	//
	VoltageLevel int64
	GsmSignle    int64
	Engine       bool
	Acc          bool
	Gps          bool
	//
	Latitude  float64
	Longitude float64
	Speed     int64
	Course    int64
}
