package main

import (
	"fmt"
	"mtrack/device_server/decoder"
	"mtrack/device_server/domain"
	"mtrack/device_server/types"
	"time"
)

func main() {

	const (
		CONN_TYPE = "tcp"
		CONN_URI  = ":1885"

		MONGODB_URI = "mongodb://root:@localhost:27017/?authMechanism=DEFAULT&authSource=admin"
	)

	channel := make(chan types.Message, 10000)

	// var db any
	store := domain.ConnectToDatabase(MONGODB_URI, "mtrucker")

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
			err = store.Vehicles.UpdateCurrentStatus(msg.IMEI,
				msg.Engine,
				msg.GsmSignle,
				msg.VoltageLevel,
				msg.Acc,
				msg.Gps,
			)
		case "LOCATION":
			// GET VEHICLE

			var e types.Vehicle
			e, err = store.Vehicles.GetVehicle(msg.IMEI)

			if err == nil {
				crnTime := time.Now().Unix()
				lastID := e.LastTripID

				tpEvent := types.TripEvent{
					Latitude:  msg.Latitude,
					Longitude: msg.Longitude,
					Speed:     msg.Speed,
					Course:    msg.Course,
				}

				if (crnTime - e.LastGpsUpdate) < 1800 {
					fmt.Println("[*] APPEND EVENT TO TRIP: ", lastID)
					err = store.Trips.AddEvent(lastID, crnTime, tpEvent)
				} else {

					fmt.Println("[*] NEW TRIP")
					lastID, err = store.Trips.CreateTrip(e.IMEI, crnTime, tpEvent)
				}

				if err == nil {
					err = store.Vehicles.UpdateLocationInfo(e.IMEI, lastID, crnTime)
				}
			}

		default:
			fmt.Println(msg.ErrorData)
		}

		if err != nil {
			fmt.Println("DB ERROR: ", err.Error())
		}
	}
}

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
