package types

type Vehicle struct {
	UID  string
	IMEI string
	Name string

	Device bool

	Engine       bool
	GsmSignle    int
	VoltageLevel int
	AccStatus    bool
	GpsStatus    bool

	Trips []string
}
