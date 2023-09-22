package types

type Message struct {
	Event string
	Data  string
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
