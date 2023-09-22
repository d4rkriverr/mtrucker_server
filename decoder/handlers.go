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

// func getVoltageLevel(v int64) string {
// 	var result = "no power (shutting down)"
// 	switch v {
// 	case 1:
// 		result = "extremely low battery"
// 	case 2:
// 		result = "very low battery (low battery alarm)"
// 	case 3:
// 		result = "low battery (can be used normally)"
// 	case 4:
// 		result = "medium"
// 	case 5:
// 		result = "high"
// 	case 6:
// 		result = "very high"
// 	default:
// 		result = "no power (shutting down)"
// 	}
// 	return result
// }

// func getGSMSingle(v int64) string {
// 	var result string
// 	switch v {
// 	case 1:
// 		result = "extremely weak signal"
// 	case 2:
// 		result = "very weak signal"
// 	case 3:
// 		result = "good signal"
// 	case 4:
// 		result = "strong signal"
// 	default:
// 		result = "no signal"
// 	}
// 	return result
// }
