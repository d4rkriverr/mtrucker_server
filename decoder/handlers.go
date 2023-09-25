package decoder

import (
	"fmt"
	"strconv"
)

type HeartBitInfo struct {
	// TerminalInfo string
	VoltageLevel int64
	GSMSignal    int64
	Engine       bool
	Acc          bool
	Gps          bool
}

func ParseHeartBit(s string) HeartBitInfo {
	var res HeartBitInfo
	// 0406030002

	x, _ := strconv.ParseInt(s[0:2], 16, 32)
	res.Engine, res.Gps, res.Acc = getTerminalInfo(fmt.Sprintf("%08s", strconv.FormatInt(x, 2)))

	res.VoltageLevel, _ = strconv.ParseInt(s[2:4], 10, 32)
	res.GSMSignal, _ = strconv.ParseInt(s[4:6], 10, 32)

	return res
}

func getTerminalInfo(s string) (bool, bool, bool) {
	return s[0] == 49, s[1] == 49, s[6] == 49
}

// ---- //
type GpsLocationInfo struct {
	Latitude  float64
	Longitude float64
	Speed     int64
	Course    int64
}

func ParseGpsLocation(s string) GpsLocationInfo {
	var result GpsLocationInfo

	data := gps_StringToArr(s)
	result.Latitude = calLatLong(data[2])
	result.Longitude = calLatLong(data[3])
	result.Speed, _ = strconv.ParseInt(data[4], 16, 64)
	// data.RealTimeGps, data.GpsPositioned, data.EastLongitude, data.NorthLatitude = calGpsStatus(s[5])
	result.Course = calGpsCourse(data[5])

	return result
}

func gps_StringToArr(s string) []string {
	var sarr []string
	sarr = append(sarr, s[:12])   // [0] DATETIME
	sarr = append(sarr, s[12:14]) // [1] Quantity of GPS information satellites
	sarr = append(sarr, s[14:22]) // [2] Latitude
	sarr = append(sarr, s[22:30]) // [3] Longitude
	sarr = append(sarr, s[30:32]) // [4] speed
	sarr = append(sarr, s[32:36]) // [5] Course / Status
	return sarr
}

func calLatLong(s string) float64 {
	r, _ := strconv.ParseInt(s, 16, 64)
	return float64(r) / float64(1800000)
}

func calGpsCourse(s string) int64 {
	b, _ := strconv.ParseInt(s[:2], 16, 32)
	c, _ := strconv.ParseInt(s[2:], 16, 32)
	sBin := fmt.Sprintf("%08b%08b", b, c)
	course, _ := strconv.ParseInt(sBin[6:], 2, 64)
	return course
}
