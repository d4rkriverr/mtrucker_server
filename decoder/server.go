package decoder

import (
	"encoding/hex"
	"fmt"
	"mtrack/device_server/types"
	"net"
	"os"

	"github.com/sigurn/crc16"
)

type MsgData struct {
	Protocol     string
	Value        string
	SerialNumber string
}

func START_TPC_SERVER(cnType, cnUri string, ch chan types.Message) {

	l, err := net.Listen(cnType, cnUri)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + cnUri)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn, ch)
	}
}

func handleRequest(cn net.Conn, ch chan types.Message) {
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	var imei string
	for {

		bufLength, err := cn.Read(buf)
		if err != nil {
			break
		}

		hexStr := hex.EncodeToString(buf[:bufLength])
		msg, err := getMessageValue(hexStr)
		if err != nil {
			ch <- types.Message{Event: "ERROR", Data: "[E1X01] unknown hex header: " + hexStr}
			continue
		}

		switch msg.Protocol {
		case "01":
			imei = msg.Value
			cn.Write(responseMessage(msg, ""))
			ch <- types.Message{Event: "LOGIN", IMEI: imei}
		case "13":
			val := ParseHeartBit(msg.Value)
			cn.Write(responseMessage(msg, ""))
			ch <- types.Message{
				Event:        "STATUS",
				IMEI:         imei,
				GsmSignle:    val.GSMSignal,
				VoltageLevel: val.VoltageLevel,
				Engine:       val.Engine,
				Acc:          val.Acc,
				Gps:          val.Gps,
			}
		default:
			ch <- types.Message{Event: "ERROR", Data: "[E1X02] unknown operation: " + hexStr}
			continue
		}
		// fmt.Println(operation, value)
	}

	if imei != "" {
		ch <- types.Message{Event: "DISCONNECT", IMEI: imei}
	}

	cn.Close()
}

func getMessageValue(s string) (MsgData, error) {
	var res MsgData

	if s[:4] == "7878" {
		res.Protocol = s[6:8]
		res.Value = s[8 : len(s)-12]
	} else if s[:4] == "7979" {
		res.Protocol = s[8:10]
		res.Value = s[10 : len(s)-12]
	} else {
		return res, fmt.Errorf("UNKNOWN REQUEST")
	}
	res.SerialNumber = s[(len(s) - 12):(len(s) - 8)]
	return res, nil
}

func responseMessage(m MsgData, r string) []byte {
	val := m.Protocol + r + m.SerialNumber
	str := fmt.Sprintf("%02x%v", (len(val)/2)+2, val)

	hm, _ := hex.DecodeString(str)
	crc := crc16.Checksum(hm, crc16.MakeTable(crc16.CRC16_X_25))

	hx := fmt.Sprintf("7878%v%x0d0a", str, crc)
	rr, _ := hex.DecodeString(hx)
	return rr
}
