package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Vehicle struct {
	UID  primitive.ObjectID `bson:"_id"`
	IMEI string
	Name string

	Device bool

	Engine       bool
	GsmSignle    int
	VoltageLevel int
	AccStatus    bool
	GpsStatus    bool

	LastTripID    string
	LastGpsUpdate int64
}

type Trip struct {
	ID         primitive.ObjectID `bson:"_id"`
	IMEI       string
	CreatedAt  int64
	LastUpdate int64
	Events     []TripEvent
}
type TripEvent struct {
	Latitude  float64
	Longitude float64
	Speed     int64
	Course    int64
}
