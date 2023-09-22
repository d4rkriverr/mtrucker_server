package main

import (
	"fmt"
	"mtrack/device_server/decoder"
	"mtrack/device_server/domain"
	"mtrack/device_server/types"
)

func main() {

	const (
		CONN_TYPE = "tcp"
		CONN_URI  = ":1885"
	)

	channel := make(chan types.Message, 10000)

	// var db any
	store := domain.ConnectToDatabase("mongodb://localhost:27017", "mtrucker")

	// start tcp server
	go decoder.START_TPC_SERVER(CONN_TYPE, CONN_URI, channel)

	// listen to devices messages // LOGIN|STATUS|GPS|ERROR

	var err error
	for {
		msg := <-channel

		switch msg.Event {
		case "LOGIN":
			err = store.Vehicles.UpdateDeviceStatus(msg.IMEI, true) // DEVICE ONLINE
		case "DISCONNECT":
			err = store.Vehicles.UpdateDeviceStatus(msg.IMEI, false) // DEVICE OFFLINE
		case "STATUS":
			fmt.Printf("%+v \n", msg)
			err = store.Vehicles.UpdateCurrentStatus(msg.IMEI,
				msg.Engine,
				msg.GsmSignle,
				msg.VoltageLevel,
				msg.Acc,
				msg.Gps,
			)
		case "GPSINFO":
		// 	notifyApi = true
		default:
			fmt.Println(msg.Data)
		}
		if err != nil {
			fmt.Println("DB ERROR: ", err.Error())
		}
	}

	// s := "" // the hex message

	// msg := decryptMessage(s) // decrypt messages

}

// func decryptMessage(s string) Message {
// 	return Message{}
// }

type Vehicle struct {
	ID   string
	IMEI string

	DeviceStatus bool // ONLINE - OFFLINE
	VoltageLevel string
	GsmSignle    string
}

type Trip struct {
	ID string

	StartAt string
	EndAt   string

	Date string

	Events string
}
